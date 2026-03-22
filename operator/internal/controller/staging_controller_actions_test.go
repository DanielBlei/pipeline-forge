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
	"testing"

	corev1alpha1 "github.com/DanielBlei/pipeline-forge/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInitializeStatus(t *testing.T) {
	ready := corev1alpha1.ObjConditionReady
	initiating := corev1alpha1.ObjConditionInitiating

	tests := []struct {
		name                   string
		staging                *corev1alpha1.Staging
		wantFirstInit          bool
		wantStatus             *corev1alpha1.ObjCondition
		wantObservedGeneration int64
	}{
		{
			name: "NilStatusSetsInitiating",
			staging: &corev1alpha1.Staging{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
			},
			wantFirstInit:          true,
			wantStatus:             &initiating,
			wantObservedGeneration: 1,
		},
		{
			name: "ExistingStatusPreserved",
			staging: &corev1alpha1.Staging{
				ObjectMeta: metav1.ObjectMeta{Generation: 3},
				Status:     corev1alpha1.StagingStatus{Status: &ready},
			},
			wantFirstInit:          false,
			wantStatus:             &ready,
			wantObservedGeneration: 3,
		},
		{
			name: "GenerationAlwaysUpdated",
			staging: &corev1alpha1.Staging{
				ObjectMeta: metav1.ObjectMeta{Generation: 5},
				Status: corev1alpha1.StagingStatus{
					Status:             &initiating,
					ObservedGeneration: 1,
				},
			},
			wantFirstInit:          false,
			wantStatus:             &initiating,
			wantObservedGeneration: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &StagingReconciler{}
			firstInit := r.initializeStatus(tt.staging)

			if firstInit != tt.wantFirstInit {
				t.Errorf("expected firstInit %v, got %v", tt.wantFirstInit, firstInit)
			}
			if tt.staging.Status.Status == nil {
				t.Fatal("expected status to be set, got nil")
			}
			if *tt.staging.Status.Status != *tt.wantStatus {
				t.Errorf("expected status %q, got %q", *tt.wantStatus, *tt.staging.Status.Status)
			}
			if tt.staging.Status.ObservedGeneration != tt.wantObservedGeneration {
				t.Errorf(
					"expected ObservedGeneration %d, got %d",
					tt.wantObservedGeneration,
					tt.staging.Status.ObservedGeneration,
				)
			}
		})
	}
}
