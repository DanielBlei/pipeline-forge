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
	"strings"
	"testing"
	"time"

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := corev1alpha1.AddToScheme(s); err != nil {
		panic("failed to add corev1alpha1 to scheme: " + err.Error())
	}
	if err := batchv1.AddToScheme(s); err != nil {
		panic("failed to add batchv1 to scheme: " + err.Error())
	}
	return s
}

func newStaging(
	name string,
	mode corev1alpha1.IngestMode,
	ingestType corev1alpha1.IngestType,
	ingestName, ingestNamespace string,
) *corev1alpha1.Staging {
	return &corev1alpha1.Staging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1alpha1.StagingSpec{
			Ingest: corev1alpha1.IngestSpec{
				Mode:      mode,
				Type:      ingestType,
				Name:      ingestName,
				Namespace: ingestNamespace,
				Schedule:  "*/5 * * * *",
				Image:     "busybox:latest",
				Command:   []string{"echo", "test"},
			},
			Transform: corev1alpha1.TransformSpec{
				Project: "test-project",
				Target:  "dev",
				Image:   "dbt:latest",
			},
		},
	}
}

func TestValidateIngestion(t *testing.T) {
	tests := []struct {
		name           string
		staging        *corev1alpha1.Staging
		existingObjs   []client.Object
		wantStatus     corev1alpha1.ObjCondition
		wantMsgContain string
		wantErr        bool
	}{
		{
			name: "ReferenceCronjobExists",
			staging: newStaging(
				"stg-ref",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				&batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{Name: "my-cj", Namespace: "default"},
				},
			},
			wantStatus: corev1alpha1.ObjConditionReady,
			wantErr:    false,
		},
		{
			name: "ReferenceCronjobMissing",
			staging: newStaging(
				"stg-missing",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"no-such-cj", "default",
			),
			existingObjs:   nil,
			wantStatus:     corev1alpha1.ObjConditionFailed,
			wantMsgContain: "not found",
			wantErr:        true,
		},
		{
			name: "BootstrapJobNotImplemented",
			staging: newStaging(
				"stg-boot-job",
				corev1alpha1.IngestModeBootstrap, corev1alpha1.IngestTypeJob,
				"boot-job", "default",
			),
			existingObjs:   nil,
			wantStatus:     corev1alpha1.ObjConditionFailed,
			wantMsgContain: "not supported",
			wantErr:        true,
		},
		{
			name: "BootstrapTriggerNotImplemented",
			staging: newStaging(
				"stg-boot-trg",
				corev1alpha1.IngestModeBootstrap, corev1alpha1.IngestTypeTrigger,
				"boot-trg", "default",
			),
			existingObjs:   nil,
			wantStatus:     corev1alpha1.ObjConditionFailed,
			wantMsgContain: "not supported",
			wantErr:        true,
		},
		{
			name: "NamespaceDefaultsToStagingNamespace",
			staging: newStaging(
				"stg-ns-fallback",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"ns-cj", "", // empty namespace, should fall back to staging's namespace
			),
			existingObjs: []client.Object{
				&batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{Name: "ns-cj", Namespace: "default"},
				},
			},
			wantStatus: corev1alpha1.ObjConditionReady,
			wantErr:    false,
		},
		{
			name: "BootstrapCronjobAlreadyExists",
			staging: newStaging(
				"stg-boot-exists",
				corev1alpha1.IngestModeBootstrap, corev1alpha1.IngestTypeCronjob,
				"existing-cj", "default",
			),
			existingObjs: []client.Object{
				&batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{Name: "existing-cj", Namespace: "default"},
				},
			},
			wantStatus: corev1alpha1.ObjConditionReady,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := newTestScheme()
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.existingObjs...).
				Build()

			r := &StagingReconciler{Client: fakeClient, Scheme: scheme}
			err := r.validateIngestion(context.Background(), tt.staging)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateIngestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.staging.Status.Ingest.Status == nil {
				t.Fatal("expected ingest status to be set, got nil")
			}
			if *tt.staging.Status.Ingest.Status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, *tt.staging.Status.Ingest.Status)
			}

			if tt.wantMsgContain != "" && !strings.Contains(tt.staging.Status.Ingest.Message, tt.wantMsgContain) {
				t.Errorf("expected message to contain %q, got %q", tt.wantMsgContain, tt.staging.Status.Ingest.Message)
			}
		})
	}
}

// cronJobUID is a fixed UID used in tests to match ownerReference filtering.
const cronJobUID = "cj-uid-12345"

// newCronJobWithUID creates a test CronJob with a fixed UID for testing Job ownership.
func newCronJobWithUID() *batchv1.CronJob {
	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cj",
			Namespace: "default",
			UID:       types.UID(cronJobUID),
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "*/5 * * * *",
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{Name: "test", Image: "busybox:latest"},
							},
						},
					},
				},
			},
		},
	}
}

// newJobOwnedByCronJob creates a Job with an ownerReference pointing to the test CronJob.
func newJobOwnedByCronJob(name string, createdAt time.Time) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(createdAt),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
					Name:       "my-cj",
					UID:        types.UID(cronJobUID),
					Controller: ptr(true),
				},
			},
		},
	}
}

func TestObserveIngestJobs(t *testing.T) {
	now := time.Now()
	completionTime := metav1.NewTime(now.Add(-1 * time.Minute))
	startTime := metav1.NewTime(now.Add(-5 * time.Minute))
	failureTime := metav1.NewTime(now.Add(-2 * time.Minute))

	tests := []struct {
		name               string
		staging            *corev1alpha1.Staging
		existingObjs       []client.Object
		wantStatus         *corev1alpha1.ObjCondition
		wantMsgContain     string
		wantLastRunJobName string
		wantSuccessful     int32
		wantFailed         int32
		wantErr            bool
	}{
		{
			name: "NoJobsStatusStaysReady",
			staging: newStaging(
				"stg-no-jobs",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				newCronJobWithUID(),
			},
			wantStatus: nil, // status unchanged (stays whatever it was before)
		},
		{
			name: "ActiveJobSetsRunning",
			staging: newStaging(
				"stg-active",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-job-1", now)
					j.Status.Active = 1
					j.Status.StartTime = &startTime
					return j
				}(),
			},
			wantStatus:         ptr(corev1alpha1.ObjConditionRunning),
			wantMsgContain:     "running",
			wantLastRunJobName: "my-cj-job-1",
		},
		{
			name: "CompletedJobSetsCompleted",
			staging: newStaging(
				"stg-completed",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-job-done", now)
					j.Status.Succeeded = 1
					j.Status.StartTime = &startTime
					j.Status.CompletionTime = &completionTime
					j.Status.Conditions = []batchv1.JobCondition{
						{
							Type:   batchv1.JobComplete,
							Status: corev1.ConditionTrue,
						},
					}
					return j
				}(),
			},
			wantStatus:         ptr(corev1alpha1.ObjectConditionCompleted),
			wantMsgContain:     "completed",
			wantLastRunJobName: "my-cj-job-done",
			wantSuccessful:     1,
		},
		{
			name: "FailedJobSetsFailed",
			staging: newStaging(
				"stg-failed",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-job-fail", now)
					j.Status.Failed = 3
					j.Status.StartTime = &startTime
					j.Status.Conditions = []batchv1.JobCondition{
						{
							Type:               batchv1.JobFailed,
							Status:             corev1.ConditionTrue,
							LastTransitionTime: failureTime,
							Message:            "BackoffLimitExceeded",
						},
					}
					return j
				}(),
			},
			wantStatus:         ptr(corev1alpha1.ObjConditionFailed),
			wantMsgContain:     "failed",
			wantLastRunJobName: "my-cj-job-fail",
			wantFailed:         1,
		},
		{
			name: "MultipleJobsPicksLatest",
			staging: newStaging(
				"stg-multi",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
				"my-cj", "default",
			),
			existingObjs: []client.Object{
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-old", now.Add(-10*time.Minute))
					j.Status.Succeeded = 1
					j.Status.Conditions = []batchv1.JobCondition{
						{Type: batchv1.JobComplete, Status: corev1.ConditionTrue},
					}
					return j
				}(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-new", now)
					j.Status.Active = 1
					j.Status.StartTime = &startTime
					return j
				}(),
			},
			wantStatus:         ptr(corev1alpha1.ObjConditionRunning),
			wantLastRunJobName: "my-cj-new",
		},
		{
			name: "SameJobTerminalSkipsObservation",
			staging: func() *corev1alpha1.Staging {
				s := newStaging(
					"stg-idempotent",
					corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob,
					"my-cj", "default",
				)
				// Simulate a previous reconcile that already processed this job
				s.Status.Ingest.LastRunJobName = "my-cj-job-done"
				s.Status.Ingest.SetInternalStatus(corev1alpha1.ObjectConditionCompleted)
				s.Status.Ingest.SuccessfulAttempts = 1
				return s
			}(),
			existingObjs: []client.Object{
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-job-done", now)
					j.Status.Succeeded = 1
					j.Status.CompletionTime = &completionTime
					j.Status.Conditions = []batchv1.JobCondition{
						{Type: batchv1.JobComplete, Status: corev1.ConditionTrue},
					}
					return j
				}(),
			},
			wantStatus:         ptr(corev1alpha1.ObjectConditionCompleted),
			wantLastRunJobName: "my-cj-job-done",
			wantSuccessful:     1, // same job terminal → short-circuit, no re-processing
		},
		{
			name: "JobTypeReturnsError",
			staging: newStaging(
				"stg-job-type",
				corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeJob,
				"my-job", "default",
			),
			existingObjs: nil,
			wantStatus:   nil,
			wantErr:      true, // observer for job type is not yet supported
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := newTestScheme()
			builder := fake.NewClientBuilder().WithScheme(scheme)
			if len(tt.existingObjs) > 0 {
				builder = builder.WithObjects(tt.existingObjs...)
			}
			fakeClient := builder.Build()

			r := &StagingReconciler{Client: fakeClient, Scheme: scheme}
			err := r.observeIngestJobs(context.Background(), tt.staging)

			if (err != nil) != tt.wantErr {
				t.Fatalf("observeIngestJobs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantStatus != nil {
				if tt.staging.Status.Ingest.Status == nil {
					t.Fatal("expected ingest status to be set, got nil")
				}
				if *tt.staging.Status.Ingest.Status != *tt.wantStatus {
					t.Errorf("expected status %q, got %q", *tt.wantStatus, *tt.staging.Status.Ingest.Status)
				}
			}

			if tt.wantMsgContain != "" &&
				!strings.Contains(strings.ToLower(tt.staging.Status.Ingest.Message), tt.wantMsgContain) {
				t.Errorf("expected message to contain %q, got %q",
					tt.wantMsgContain, tt.staging.Status.Ingest.Message)
			}

			if tt.wantLastRunJobName != "" && tt.staging.Status.Ingest.LastRunJobName != tt.wantLastRunJobName {
				t.Errorf("expected LastRunJobName %q, got %q",
					tt.wantLastRunJobName, tt.staging.Status.Ingest.LastRunJobName)
			}

			if tt.staging.Status.Ingest.SuccessfulAttempts != tt.wantSuccessful {
				t.Errorf("expected SuccessfulAttempts %d, got %d",
					tt.wantSuccessful, tt.staging.Status.Ingest.SuccessfulAttempts)
			}

			if tt.staging.Status.Ingest.FailedAttempts != tt.wantFailed {
				t.Errorf("expected FailedAttempts %d, got %d",
					tt.wantFailed, tt.staging.Status.Ingest.FailedAttempts)
			}
		})
	}
}

func TestFindLatestOwnedJob(t *testing.T) {
	now := time.Now()
	uid := types.UID(cronJobUID)
	otherUID := types.UID("other-uid")

	ownedRef := []metav1.OwnerReference{
		{UID: uid, Controller: ptr(true)},
	}

	t.Run("EmptySliceReturnsNil", func(t *testing.T) {
		result := findLatestOwnedJob(nil, uid)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("SingleOwnedJobReturnsIt", func(t *testing.T) {
		jobs := []batchv1.Job{
			{ObjectMeta: metav1.ObjectMeta{
				Name: "only", CreationTimestamp: metav1.NewTime(now),
				OwnerReferences: ownedRef,
			}},
		}
		result := findLatestOwnedJob(jobs, uid)
		if result == nil || result.Name != "only" {
			t.Errorf("expected 'only', got %v", result)
		}
	})

	t.Run("SkipsUnownedJobs", func(t *testing.T) {
		jobs := []batchv1.Job{
			{ObjectMeta: metav1.ObjectMeta{
				Name: "unowned", CreationTimestamp: metav1.NewTime(now),
				OwnerReferences: []metav1.OwnerReference{
					{UID: otherUID, Controller: ptr(true)},
				},
			}},
		}
		result := findLatestOwnedJob(jobs, uid)
		if result != nil {
			t.Errorf("expected nil, got %v", result.Name)
		}
	})

	t.Run("MultipleJobsReturnsNewestOwned", func(t *testing.T) {
		jobs := []batchv1.Job{
			{ObjectMeta: metav1.ObjectMeta{
				Name: "old", CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Minute)),
				OwnerReferences: ownedRef,
			}},
			{ObjectMeta: metav1.ObjectMeta{
				Name: "newest-unowned", CreationTimestamp: metav1.NewTime(now),
				OwnerReferences: []metav1.OwnerReference{
					{UID: otherUID, Controller: ptr(true)},
				},
			}},
			{ObjectMeta: metav1.ObjectMeta{
				Name: "newest-owned", CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Minute)),
				OwnerReferences: ownedRef,
			}},
		}
		result := findLatestOwnedJob(jobs, uid)
		if result == nil || result.Name != "newest-owned" {
			t.Errorf("expected 'newest-owned', got %v", result)
		}
	})
}

func TestPollIntervalGating(t *testing.T) {
	now := time.Now()
	completionTime := metav1.NewTime(now.Add(-1 * time.Minute))

	t.Run("SkipsWhenPollIntervalNotElapsed", func(t *testing.T) {
		staging := newStaging(
			"stg-poll", corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob, "my-cj", "default",
		)
		// 30s ago — default poll is 120s, so should skip
		staging.Status.Ingest.LastCheckedTime = metav1.Time{Time: now.Add(-30 * time.Second)}

		scheme := newTestScheme()
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
		r := &StagingReconciler{Client: fakeClient, Scheme: scheme}

		err := r.observeIngestJobs(context.Background(), staging)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		// LastCheckedTime should NOT be updated (observation was skipped)
		if time.Since(staging.Status.Ingest.LastCheckedTime.Time) < 25*time.Second {
			t.Error("LastCheckedTime should not have been updated")
		}
	})

	t.Run("RunsWhenPollIntervalElapsed", func(t *testing.T) {
		staging := newStaging(
			"stg-poll-elapsed", corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob, "my-cj", "default",
		)
		// 130s ago — default poll is 120s, so should run
		staging.Status.Ingest.LastCheckedTime = metav1.Time{Time: now.Add(-130 * time.Second)}

		scheme := newTestScheme()
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).
			WithObjects(newCronJobWithUID()).
			Build()
		r := &StagingReconciler{Client: fakeClient, Scheme: scheme}

		err := r.observeIngestJobs(context.Background(), staging)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		// LastCheckedTime should be updated
		if time.Since(staging.Status.Ingest.LastCheckedTime.Time) > 2*time.Second {
			t.Error("LastCheckedTime should have been updated to now")
		}
	})

	t.Run("RespectsCustomPollInterval", func(t *testing.T) {
		staging := newStaging(
			"stg-custom-poll", corev1alpha1.IngestModeReference, corev1alpha1.IngestTypeCronjob, "my-cj", "default",
		)
		staging.Spec.Ingest.PollIntervalSeconds = 30
		// 20s ago — custom poll is 30s, so should skip
		staging.Status.Ingest.LastCheckedTime = metav1.Time{Time: now.Add(-20 * time.Second)}

		scheme := newTestScheme()
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
		r := &StagingReconciler{Client: fakeClient, Scheme: scheme}

		err := r.observeIngestJobs(context.Background(), staging)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		// Should have skipped — 20s < 30s
		if time.Since(staging.Status.Ingest.LastCheckedTime.Time) < 15*time.Second {
			t.Error("observation should have been skipped")
		}
	})

	t.Run("SameJobRunningRechecks", func(t *testing.T) {
		staging := newStaging(
			"stg-running-recheck", corev1alpha1.IngestModeReference,
			corev1alpha1.IngestTypeCronjob, "my-cj", "default",
		)
		staging.Status.Ingest.LastRunJobName = "my-cj-job-1"
		staging.Status.Ingest.SetInternalStatus(corev1alpha1.ObjConditionRunning)

		scheme := newTestScheme()
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).
			WithObjects(
				newCronJobWithUID(),
				func() *batchv1.Job {
					j := newJobOwnedByCronJob("my-cj-job-1", now)
					j.Status.Succeeded = 1
					j.Status.CompletionTime = &completionTime
					j.Status.Conditions = []batchv1.JobCondition{
						{Type: batchv1.JobComplete, Status: corev1.ConditionTrue},
					}
					return j
				}(),
			).Build()
		r := &StagingReconciler{Client: fakeClient, Scheme: scheme}

		err := r.observeIngestJobs(context.Background(), staging)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		// Same job but was Running → should re-check and update to Completed
		if staging.Status.Ingest.Status == nil || *staging.Status.Ingest.Status != corev1alpha1.ObjectConditionCompleted {
			t.Errorf("expected Completed, got %v", staging.Status.Ingest.Status)
		}
	})
}

func ptr[T any](v T) *T {
	return &v
}
