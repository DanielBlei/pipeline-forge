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

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
			wantMsgContain: "not implemented",
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
			wantMsgContain: "not implemented",
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

			if tt.staging.Status.Ingest.LastCheckedTime.IsZero() {
				t.Error("expected LastCheckedTime to be set, got zero")
			}
		})
	}
}
