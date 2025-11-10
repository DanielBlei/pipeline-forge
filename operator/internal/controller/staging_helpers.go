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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// checkResourceExists validates that a referenced resource exists, or handles bootstrap mode
func (r *StagingReconciler) checkResourceExists(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
	obj client.Object,
) error {
	namespace := stagingObj.Spec.Ingest.Namespace
	if namespace == "" {
		namespace = stagingObj.Namespace
	}

	namespacedName := types.NamespacedName{
		Name:      stagingObj.Spec.Ingest.Name,
		Namespace: namespace,
	}

	// Use the type from the spec - it's already in the right format for error messages
	resourceType := string(stagingObj.Spec.Ingest.Type)

	if err := r.Get(ctx, namespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			if stagingObj.Spec.Ingest.Mode == corev1alpha1.IngestModeReference {
				return fmt.Errorf(
					"%s %s not found, consider using bootstrap mode to create the %s",
					resourceType,
					stagingObj.Spec.Ingest.Name,
					resourceType,
				)
			}
			if stagingObj.Spec.Ingest.Mode == corev1alpha1.IngestModeBootstrap {
				switch obj.(type) {
				case *batchv1.CronJob:
					return r.Create(ctx, components.CreateCronJob(
						stagingObj.Spec.Ingest.Name,
						stagingObj.Namespace,
						stagingObj.Spec.Ingest.Schedule,
						stagingObj.Spec.Ingest.Image,
						stagingObj.Spec.Ingest.Command,
						stagingObj.Spec.Ingest.Resources,
					))
					// TODO: add support for Job and Trigger
				}
				return err
			}
		}
		return fmt.Errorf(
			"failed to check resource %s in namespace %s",
			namespacedName.Name,
			namespacedName.Namespace,
		)
	}
	return nil
}

func (r *StagingReconciler) checkCronjobCreateIfNeeded(ctx context.Context, stagingObj *corev1alpha1.Staging) error {
	return r.checkResourceExists(ctx, stagingObj, &batchv1.CronJob{})
}

func (r *StagingReconciler) checkJobCreateIfNeeded(ctx context.Context, stagingObj *corev1alpha1.Staging) error {
	return r.checkResourceExists(ctx, stagingObj, &batchv1.Job{})
}

func (r *StagingReconciler) checkTriggerCreateIfNeeded(ctx context.Context, stagingObj *corev1alpha1.Staging) error {
	return r.checkResourceExists(ctx, stagingObj, &corev1alpha1.Trigger{})
}

func (r *StagingReconciler) validateIngestionAndUpdateStatus(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) error {
	var validationError error
	switch stagingObj.Spec.Ingest.Type {
	case corev1alpha1.IngestTypeCronjob:
		validationError = r.checkCronjobCreateIfNeeded(ctx, stagingObj)
	case corev1alpha1.IngestTypeJob:
		validationError = r.checkJobCreateIfNeeded(ctx, stagingObj)
	case corev1alpha1.IngestTypeTrigger:
		validationError = r.checkTriggerCreateIfNeeded(ctx, stagingObj)
	default:
		return fmt.Errorf("invalid ingestion type: %s", stagingObj.Spec.Ingest.Type)
	}
	if validationError == nil {
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionReady)
		stagingObj.Status.Ingest.LastCheckedTime = metav1.Time{Time: time.Now()}
		return validationError
	}
	stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionFailed)
	stagingObj.Status.Ingest.LastCheckedTime = metav1.Time{Time: time.Now()}
	stagingObj.Status.Ingest.Message = validationError.Error()
	return validationError
}
