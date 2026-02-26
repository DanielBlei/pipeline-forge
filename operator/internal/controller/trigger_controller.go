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

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	"github.com/DanielBlei/pipeline-forge/operator/internal/status"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// TriggerReconciler reconciles a Trigger object
type TriggerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=triggers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=triggers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.pipeline-forge.io,resources=triggers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Trigger object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *TriggerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx).WithValues("trigger", req.NamespacedName)
	log.Info("Reconciling Trigger")

	triggerObj := &corev1alpha1.Trigger{}
	if err := r.Get(ctx, req.NamespacedName, triggerObj); err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("Trigger object not found, stopping reconciliation")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to retrieve Trigger object")
		return ctrl.Result{}, err
	}

	// Create deep copy for status updates
	triggerDeepCopy := triggerObj.DeepCopy()

	// Initialize status if needed
	if triggerObj.Status.Status == nil {
		log.Info("New Trigger object, initializing status")
		triggerObj.Status.SetTriggerStatus(corev1alpha1.ObjConditionInitiating)
		triggerObj.Status.ObservedGeneration = triggerObj.Generation

		if err := status.UpdateStatus(ctx, r.Client, triggerObj, triggerDeepCopy); err != nil {
			log.Error(err, "Unable to update Trigger status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil // Requeue to continue processing
	}

	// Handle deletion
	if !triggerObj.DeletionTimestamp.IsZero() {
		log.Info("Trigger is being deleted")
		// TODO: Add deletion logic (cleanup resources, remove finalizers)
		return ctrl.Result{}, nil
	}

	// TODO: Add actual reconciliation logic here
	// - Check trigger conditions (BigQuery, GCS, PubSub)
	// - Evaluate trigger criteria
	// - Create/update resources as needed

	// Update observed generation
	triggerObj.Status.ObservedGeneration = triggerObj.Generation

	log.Info("Trigger reconciliation completed")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TriggerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Trigger{}).
		Named("trigger").
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.AnnotationChangedPredicate{},
			)).
		Complete(r)
}
