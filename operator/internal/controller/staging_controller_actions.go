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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// fetchStaging retrieves the Staging object from the cluster.
// Returns (nil, nil) when the object has been deleted.
func (r *StagingReconciler) fetchStaging(
	ctx context.Context,
	name types.NamespacedName,
) (*corev1alpha1.Staging, error) {
	log := logf.FromContext(ctx)

	stagingObj := &corev1alpha1.Staging{}
	if err := r.Get(ctx, name, stagingObj); err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("Staging object not found, stopping reconciliation")
			return nil, nil
		}
		log.Error(err, "Unable to retrieve Staging object")
		return nil, err
	}
	return stagingObj, nil
}

// ensureFinalizer adds the finalizer if not already present and requeues.
func (r *StagingReconciler) ensureFinalizer(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Adding finalizer")
	controllerutil.AddFinalizer(stagingObj, stagingFinalizer)
	if err := r.Update(ctx, stagingObj); err != nil {
		log.Error(err, "Failed to add finalizer")
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: time.Second}, nil
}

// handleDeletion removes the finalizer and performs cleanup when the resource is being deleted.
func (r *StagingReconciler) handleDeletion(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Staging is being deleted, running cleanup")
	if controllerutil.ContainsFinalizer(stagingObj, stagingFinalizer) {
		controllerutil.RemoveFinalizer(stagingObj, stagingFinalizer)
		if err := r.Update(ctx, stagingObj); err != nil {
			log.Error(err, "Failed to remove finalizer")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// initializeStatus sets the initial status and observed generation.
// Returns true if this was the first initialization (status was nil).
func (r *StagingReconciler) initializeStatus(stagingObj *corev1alpha1.Staging) bool {
	firstInit := stagingObj.Status.Status == nil
	if firstInit {
		stagingObj.Status.SetStagingStatus(corev1alpha1.ObjConditionInitiating)
	}
	stagingObj.Status.ObservedGeneration = stagingObj.Generation
	return firstInit
}

// patchStatusAndRequeue patches any status update and requeues.
// On success (failure == nil) it increments the "requeued" metric and requeues after defaultRequeueAfter.
// On failure it returns the error immediately so the controller-runtime backoff takes over.
func (r *StagingReconciler) patchStatusAndRequeue(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
	stagingDeepCopy *corev1alpha1.Staging,
	failure error,
) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if err := status.UpdateStatus(ctx, r.Client, stagingObj, stagingDeepCopy); err != nil {
		log.Error(err, "Unable to update Staging status")
		return ctrl.Result{}, err
	}
	if failure != nil {
		return ctrl.Result{}, failure
	}
	pfmetrics.StagingReconcileTotal.WithLabelValues(pfmetrics.ResultRequeued).Inc()
	return ctrl.Result{RequeueAfter: defaultRequeueAfter}, nil
}
