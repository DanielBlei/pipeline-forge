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
	corev1 "k8s.io/api/core/v1"
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
		return fmt.Errorf("bootstrap mode for %s is not supported", stagingObj.Spec.Ingest.Type)
	case *corev1alpha1.Trigger:
		return fmt.Errorf("bootstrap mode for %s is not supported", stagingObj.Spec.Ingest.Type)
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

	if validationError != nil {
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionFailed)
		stagingObj.Status.Ingest.Message = validationError.Error()
		return validationError
	}

	stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionReady)
	return nil
}

// observeIngestJobs dispatches to type-specific observation methods to track
// the status of Jobs spawned by the ingest resource.
// Skips observation if the poll interval has not elapsed since LastCheckedTime.
// LastCheckedTime is updated only after a successful observation.
func (r *StagingReconciler) observeIngestJobs(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) error {
	log := logf.FromContext(ctx)
	pollInterval := stagingObj.IngestPollInterval()
	if !stagingObj.Status.Ingest.LastCheckedTime.IsZero() &&
		time.Since(stagingObj.Status.Ingest.LastCheckedTime.Time) < pollInterval {
		log.V(1).Info("Poll interval not elapsed, skipping observation",
			"lastChecked", stagingObj.Status.Ingest.LastCheckedTime.Time,
			"interval", pollInterval,
		)
		return nil
	}

	log.Info("Observing ingest jobs",
		"type", stagingObj.Spec.Ingest.Type,
		"name", stagingObj.Spec.Ingest.Name,
	)

	var err error
	switch stagingObj.Spec.Ingest.Type {
	case corev1alpha1.IngestTypeCronjob:
		err = r.observeCronJobJobs(ctx, stagingObj)
	case corev1alpha1.IngestTypeJob:
		err = fmt.Errorf("observer for %s is not supported", corev1alpha1.IngestTypeJob)
	case corev1alpha1.IngestTypeTrigger:
		err = fmt.Errorf("observer for %s is not supported", corev1alpha1.IngestTypeTrigger)
	default:
		err = fmt.Errorf("invalid obj type %s", stagingObj.Spec.Ingest.Type)
	}
	if err == nil {
		stagingObj.Status.Ingest.LastCheckedTime = metav1.Time{Time: time.Now()}
	}
	return err
}

// observeCronJobJobs lists Jobs owned by the ingest CronJob and updates the ingest status
// based on the latest Job's state.
func (r *StagingReconciler) observeCronJobJobs(
	ctx context.Context,
	stagingObj *corev1alpha1.Staging,
) error {
	log := logf.FromContext(ctx)
	namespace := stagingObj.IngestNamespace()

	cronJob := &batchv1.CronJob{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      stagingObj.Spec.Ingest.Name,
		Namespace: namespace,
	}, cronJob); err != nil {
		return fmt.Errorf("failed to get CronJob %s/%s: %w", namespace, stagingObj.Spec.Ingest.Name, err)
	}

	jobList := &batchv1.JobList{}
	if err := r.List(ctx, jobList, client.InNamespace(namespace)); err != nil {
		return fmt.Errorf("failed to list jobs in %s: %w", namespace, err)
	}

	latestJob := findLatestOwnedJob(jobList.Items, cronJob.UID)
	if latestJob == nil {
		log.V(1).Info("No jobs found for CronJob, waiting for first schedule", "cronjob", cronJob.Name)
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionPending)
		stagingObj.Status.Ingest.Message = fmt.Sprintf("Waiting for CronJob %s to schedule its first Job", cronJob.Name)
		return nil
	}

	if latestJob.Name == stagingObj.Status.Ingest.LastRunJobName {
		currentStatus := stagingObj.Status.Ingest.Status
		if currentStatus != nil &&
			(*currentStatus == corev1alpha1.ObjectConditionCompleted ||
				*currentStatus == corev1alpha1.ObjConditionFailed) {
			log.V(1).Info("Latest job unchanged, skipping", "job", latestJob.Name)
			return nil
		}
	}
	updateIngestStatusFromJob(stagingObj, latestJob)
	return nil
}

// updateIngestStatusFromJob updates the ingest InternalStatus based on the given Job's status.
// Counter increments are gated on LastRunJobName changing to ensure idempotency across re-reconciles.
func updateIngestStatusFromJob(stagingObj *corev1alpha1.Staging, job *batchv1.Job) {
	isNewJob := stagingObj.Status.Ingest.LastRunJobName != job.Name
	stagingObj.Status.Ingest.LastRunJobName = job.Name

	if job.Status.StartTime != nil {
		stagingObj.Status.Ingest.LastAttemptTime = job.Status.StartTime
	}

	for _, cond := range job.Status.Conditions {
		if cond.Status != corev1.ConditionTrue {
			continue
		}

		switch cond.Type {
		case batchv1.JobComplete:
			stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjectConditionCompleted)
			stagingObj.Status.Ingest.Message = fmt.Sprintf("Job %s completed", job.Name)
			if job.Status.CompletionTime != nil {
				stagingObj.Status.Ingest.LastCompletedTime = *job.Status.CompletionTime
			}
			if isNewJob {
				stagingObj.Status.Ingest.SuccessfulAttempts++
			}
			return

		case batchv1.JobFailed:
			stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionFailed)
			stagingObj.Status.Ingest.Message = fmt.Sprintf("Job %s failed: %s", job.Name, cond.Message)
			stagingObj.Status.Ingest.LastFailureTime = &cond.LastTransitionTime
			if isNewJob {
				stagingObj.Status.Ingest.FailedAttempts++
			}
			return
		}
	}

	if job.Status.Active > 0 {
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionRunning)
		stagingObj.Status.Ingest.Message = fmt.Sprintf("Job %s is running", job.Name)
	} else {
		stagingObj.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionPending)
		stagingObj.Status.Ingest.Message = fmt.Sprintf("Job %s is pending", job.Name)
	}
}

// findLatestOwnedJob returns the most recently created Job owned by the given controller UID,
// or nil if no matching Jobs exist. Single linear scan over the slice.
func findLatestOwnedJob(jobs []batchv1.Job, ownerUID types.UID) *batchv1.Job {
	var latest *batchv1.Job
	for i := range jobs {
		if !isOwnedByController(&jobs[i], ownerUID) {
			continue
		}
		if latest == nil || jobs[i].CreationTimestamp.After(latest.CreationTimestamp.Time) {
			latest = &jobs[i]
		}
	}
	return latest
}

func isOwnedByController(job *batchv1.Job, uid types.UID) bool {
	for _, ref := range job.OwnerReferences {
		if ref.UID == uid && ref.Controller != nil && *ref.Controller {
			return true
		}
	}
	return false
}
