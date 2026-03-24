/*
Copyright 2025 Daniel Blei.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	pfmetrics "github.com/DanielBlei/pipeline-forge/operator/internal/metrics"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	stagingFinalizer    = "staging.core.pipeline-forge.io/finalizer"
	defaultRequeueAfter = 60 * time.Second
)

// StagingReconciler reconciles a Staging object
type StagingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=stagings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=stagings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=stagings/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;delete

func (r *StagingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx).WithValues("staging", req.NamespacedName)
	log.Info("Reconciling Staging")

	stagingObj, err := r.fetchStaging(ctx, req.NamespacedName)
	if stagingObj == nil || err != nil {
		return ctrl.Result{}, err
	}

	if !stagingObj.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, stagingObj)
	}

	if !controllerutil.ContainsFinalizer(stagingObj, stagingFinalizer) {
		return r.ensureFinalizer(ctx, stagingObj)
	}

	stagingDeepCopy := stagingObj.DeepCopy()

	if r.initializeStatus(stagingObj) {
		pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultInitialized).Inc()
	}

	if err := r.validateIngestion(ctx, stagingObj); err != nil {
		log.Error(err, "Error checking ingestion")
		pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultValidationFailed).Inc()
		return r.patchStatusAndRequeue(ctx, stagingObj, stagingDeepCopy, err)
	}

	if err := r.observeIngestJobs(ctx, stagingObj); err != nil {
		log.Error(err, "Error observing ingest jobs")
		return r.patchStatusAndRequeue(ctx, stagingObj, stagingDeepCopy, err)
	}

	// TODO: validate transformation and perform transformation

	return r.patchStatusAndRequeue(ctx, stagingObj, stagingDeepCopy, nil)
}

// SetupWithManager sets up the controller with the Manager.
func (r *StagingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := pfmetrics.Register(); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Staging{}).
		Named("staging").
		// Only reconcile when spec changes (generation bump) or annotations change.
		// This prevents unnecessary reconciliations on status-only updates.
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.AnnotationChangedPredicate{},
			)).
		Complete(r)
}
