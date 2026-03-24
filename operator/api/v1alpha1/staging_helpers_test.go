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

package v1alpha1

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIngestNamespace(t *testing.T) {
	t.Run("ReturnsSpecNamespace", func(t *testing.T) {
		s := &Staging{
			ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
			Spec:       StagingSpec{Ingest: IngestSpec{Namespace: "custom"}},
		}
		if got := s.IngestNamespace(); got != "custom" {
			t.Errorf("expected 'custom', got %q", got)
		}
	})

	t.Run("FallsBackToStagingNamespace", func(t *testing.T) {
		s := &Staging{
			ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
			Spec:       StagingSpec{Ingest: IngestSpec{}},
		}
		if got := s.IngestNamespace(); got != "default" {
			t.Errorf("expected 'default', got %q", got)
		}
	})
}

func TestIngestPollInterval(t *testing.T) {
	t.Run("DefaultsTo120s", func(t *testing.T) {
		s := &Staging{}
		if got := s.IngestPollInterval(); got != 120*time.Second {
			t.Errorf("expected 120s, got %v", got)
		}
	})

	t.Run("UsesSpecValue", func(t *testing.T) {
		s := &Staging{}
		s.Spec.Ingest.PollIntervalSeconds = 60
		if got := s.IngestPollInterval(); got != 60*time.Second {
			t.Errorf("expected 60s, got %v", got)
		}
	})
}

func TestIngestSuspended(t *testing.T) {
	t.Run("ReturnsTrueWhenSuspended", func(t *testing.T) {
		s := &Staging{Spec: StagingSpec{Ingest: IngestSpec{Suspend: true}}}
		if !s.IngestSuspended() {
			t.Error("expected true")
		}
	})

	t.Run("ReturnsFalseWhenNotSuspended", func(t *testing.T) {
		s := &Staging{}
		if s.IngestSuspended() {
			t.Error("expected false")
		}
	})
}

func TestTransformSuspended(t *testing.T) {
	t.Run("ReturnsTrueWhenSuspended", func(t *testing.T) {
		s := &Staging{
			Spec: StagingSpec{Transform: TransformSpec{Suspend: true}},
		}
		if !s.TransformSuspended() {
			t.Error("expected true")
		}
	})

	t.Run("ReturnsFalseWhenNotSuspended", func(t *testing.T) {
		s := &Staging{}
		if s.TransformSuspended() {
			t.Error("expected false")
		}
	})
}
