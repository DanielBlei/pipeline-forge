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
	"fmt"
	"time"

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	"github.com/DanielBlei/pipeline-forge/operator/internal/components"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// checkResourceReference verifies that a referenced resource exists in the cluster.
func (r *StagingReconciler) checkResourceReference(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
	obj client.Object,
) error {
	namespacedName := types.NamespacedName{
		Name:      stagingObj.Spec.Ingest.Name,
		Namespace: stagingObj.IngestNamespace(),
	}

	if err := r.Get(ctx, namespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf(
				"referenced %s %s/%s not found",
				stagingObj.Spec.Ingest.Type,
				namespacedName.Namespace,
				namespacedName.Name,
			)
		}
		return fmt.Errorf(
			"failed to check resource %s/%s: %w",
			namespacedName.Namespace,
			namespacedName.Name,
			err,
		)
	}
	return nil
}

// bootstrapResource creates the ingest resource if it does not already exist.
// TODO: move to a createOrUpdate resource approach to enforce desired vs actual state
func (r *StagingReconciler) bootstrapResource(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
	obj client.Object,
) error {
	namespacedName := types.NamespacedName{
		Name:      stagingObj.Spec.Ingest.Name,
		Namespace: stagingObj.IngestNamespace(),
	}

	if err := r.Get(ctx, namespacedName, obj); err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf(
				"failed to check resource %s/%s: %w",
				namespacedName.Namespace,
				namespacedName.Name,
				err,
			)
		}
		return r.createResource(ctx, stagingObj, obj)
	}
	// Resource already exists, nothing to do.
	return nil
}

// createResource creates the appropriate Kubernetes resource for bootstrap mode.
func (r *StagingReconciler) createResource(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
	obj client.Object,
) error {
	switch obj.(type) {
	case *batchv1.CronJob:
		cronJob := components.CreateCronJob(
			stagingObj.Spec.Ingest.Name,
			stagingObj.IngestNamespace(),
			stagingObj.Spec.Ingest.Schedule,
			stagingObj.Spec.Ingest.Image,
			stagingObj.Spec.Ingest.Command,
			stagingObj.Spec.Ingest.Resources,
		)
		if err := ctrl.SetControllerReference(stagingObj, cronJob, r.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference on CronJob: %w", err)
		}
		return r.Create(ctx, cronJob)
	case *batchv1.Job:
		return fmt.Errorf("bootstrap mode for %s is not implemented", stagingObj.Spec.Ingest.Type)
	case *corev1alpha1.Trigger:
		return fmt.Errorf("bootstrap mode for %s is not implemented", stagingObj.Spec.Ingest.Type)
	default:
		return fmt.Errorf("bootstrap mode is not supported for %s", stagingObj.Spec.Ingest.Type)
	}
}

// resolveIngestResource dispatches to the correct check or bootstrap handler based on mode and type.
func (r *StagingReconciler) resolveIngestResource(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) error {
	type resolver struct {
		obj client.Object
	}

	resolvers := map[corev1alpha1.IngestType]resolver{
		corev1alpha1.IngestTypeCronjob: {obj: &batchv1.CronJob{}},
		corev1alpha1.IngestTypeJob:     {obj: &batchv1.Job{}},
		corev1alpha1.IngestTypeTrigger: {obj: &corev1alpha1.Trigger{}},
	}

	res, ok := resolvers[stagingObj.Spec.Ingest.Type]
	if !ok {
		return fmt.Errorf("invalid ingestion type: %s", stagingObj.Spec.Ingest.Type)
	}

	switch stagingObj.Spec.Ingest.Mode {
	case corev1alpha1.IngestModeReference:
		return r.checkResourceReference(ctx, stagingObj, res.obj)
	case corev1alpha1.IngestModeBootstrap:
		return r.bootstrapResource(ctx, stagingObj, res.obj)
	default:
		return fmt.Errorf("invalid ingest mode: %s", stagingObj.Spec.Ingest.Mode)
	}
}

func (r *StagingReconciler) validateIngestion(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) error {
	log := logf.FromContext(ctx)
	log.Info("Validating ingestion",
		"mode", stagingObj.Spec.Ingest.Mode,
		"type", stagingObj.Spec.Ingest.Type,
		"name", stagingObj.Spec.Ingest.Name,
	)

	validationError := r.resolveIngestResource(ctx, stagingObj)

	stagingObj.Status.Ingest.LastCheckedTime = metav1.Time{Time: time.Now()}
	if validationError != nil {
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionFailed)
		stagingObj.Status.Ingest.Message = validationError.Error()
		return validationError
	}

	stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionReady)
	return nil
}
