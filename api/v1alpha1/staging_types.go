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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJobReference defines a reference to an existing CronJob resource.
type CronJobReference struct {
	// Kind of the referenced object (usually "CronJob")
	// +kubebuilder:validation:Enum=CronJob
	Kind string `json:"kind"`

	// Namespace where the referenced CronJob exists.
	// If omitted, defaults to the namespace of the Staging resource.
	Namespace string `json:"namespace,omitempty"`

	// Name of the referenced CronJob
	Name string `json:"name"`
}

// TriggerReference specifies the location of a Trigger resource used to activate this stage.
type TriggerReference struct {
	// Namespace where the Trigger resource is defined.
	// If omitted, defaults to the namespace of the Staging resource.
	Namespace string `json:"namespace,omitempty"`

	// Name of the Trigger resource (required)
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// IngestSpec defines how this stage is initiated by an upstream job or event.
type IngestSpec struct {
	// Mode determines how this stage is triggered: by a CronJob reference or a Trigger CRD
	// +kubebuilder:validation:Enum=reference;trigger
	// +kubebuilder:validation:Required
	Mode string `json:"mode"`

	// Owner indicates who owns or maintains this ingestion stage (optional metadata)
	Owner string `json:"owner,omitempty"`

	// Reference to an existing CronJob resource (used if mode=reference)
	// +optional
	Reference *CronJobReference `json:"cronjob_ref,omitempty"`

	// Reference to a Trigger resource (used if mode=trigger)
	// +optional
	Trigger *TriggerReference `json:"trigger,omitempty"`
}

// TransformSpec defines the dbt transformation to execute after ingestion.
type TransformSpec struct {
	// Name of the dbt project this stage executes
	// +kubebuilder:validation:Required
	Project string `json:"project"`

	// dbt target profile to use (e.g. "dev", "prod")
	// +kubebuilder:validation:Required
	Target string `json:"target"`

	// Optional owner metadata for the transform stage
	Owner string `json:"owner,omitempty"`

	// List of dbt models to execute
	// +kubebuilder:validation:MinItems=1
	Models []string `json:"models"`

	// Whether to run dbt with --full-refresh
	FullRefresh bool `json:"full_refresh,omitempty"`
}

// StagingSpec defines the desired state of Staging
type StagingSpec struct {
	// Optional human-readable description of this staging step
	Description string `json:"description,omitempty"`

	// Ingest defines how the pipeline is initiated (e.g., via CronJob or Trigger)
	// +kubebuilder:validation:Required
	Ingest IngestSpec `json:"ingest"`

	// Transform defines the dbt-based transformation logic
	// +kubebuilder:validation:Required
	Transform TransformSpec `json:"transform"`
}

// StagingStatus defines the observed state of Staging.
type StagingStatus struct {
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Staging is the Schema for the stagings API
type Staging struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of Staging
	Spec StagingSpec `json:"spec"`

	// Status defines the observed state of Staging
	Status StagingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StagingList contains a list of Staging
type StagingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Staging `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Staging{}, &StagingList{})
}
