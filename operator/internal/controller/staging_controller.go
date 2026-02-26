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
	"github.com/DanielBlei/pipeline-forge/operator/internal/status"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	stagingObj := &corev1alpha1.Staging{}
	if err := r.Get(ctx, req.NamespacedName, stagingObj); err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("Staging object not found, stopping reconciliation")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to retrieve Staging object")
		return ctrl.Result{}, err
	}

	// Handle deletion first
	if !stagingObj.DeletionTimestamp.IsZero() {
		log.Info("Staging is being deleted")
		// TODO: Add deletion logic (cleanup resources, remove finalizers)
		return ctrl.Result{}, nil
	}

	// Create deep copy  prior to object checks and updates
	stagingDeepCopy := stagingObj.DeepCopy()

	// Initialize status if needed (first time reconciliation)
	if stagingObj.Status.Status == nil {
		log.Info("New Staging object, initializing status")
		stagingObj.Status.SetStagingStatus(corev1alpha1.ObjConditionInitiating)
		pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultInitialized).Inc()
	}

	// Update observed generation
	stagingObj.Status.ObservedGeneration = stagingObj.Generation

	// Check Ingestion (work in progress)
	if err := r.validateIngestionAndUpdateStatus(ctx, stagingObj); err != nil {
		log.Error(err, "Error checking ingestion")
		if updateError := status.UpdateStatus(ctx, r.Client, stagingObj, stagingDeepCopy); updateError != nil {
			log.Error(updateError, "Unable to update Staging status")
			return ctrl.Result{}, updateError
		}
		pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultValidationFailed).Inc()
		return ctrl.Result{}, err
	}

	if err := status.UpdateStatus(ctx, r.Client, stagingObj, stagingDeepCopy); err != nil {
		log.Error(err, "Unable to update Staging status")
		return ctrl.Result{}, err
	}

	pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultRequeued).Inc()
	return ctrl.Result{RequeueAfter: defaultRequeueAfter}, nil
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
