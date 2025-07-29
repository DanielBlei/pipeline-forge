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

// IngestSpec defines how this stage is initiated by an upstream job or event.
type IngestSpec struct {
	// Mode determines how this stage is triggered: by a CronJob reference or a Trigger CRD
	// +kubebuilder:validation:Enum=cronjob;job;trigger
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// Reference to an existing CronJob resource (used if mode=reference)
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Namespace where the resource is defined.
	// If omitted, defaults to the namespace of the Staging resource.
	Namespace string `json:"namespace,omitempty"`

	// Owner indicates who owns or maintains this ingestion stage (optional metadata)
	// +optional
	Owner string `json:"owner,omitempty"`
}

// TransformSpec defines the dbt transformation to execute after ingestion.
type TransformSpec struct {
	// Name of the dbt project this stage executes
	// +kubebuilder:validation:Required
	Project string `json:"project"`

	// dbt target profile to use (e.g. "dev", "prod")
	// +kubebuilder:validation:Required
	Target string `json:"target"`

	// Image to run the Transform
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Optional owner metadata for the transform stage
	Owner string `json:"owner,omitempty"`

	// List of dbt models to execute
	// +kubebuilder:validation:MinItems=1
	Models []string `json:"models"`

	// Whether to run dbt with --full-refresh
	// +optional
	FullRefresh bool `json:"full_refresh,omitempty"`
}

// StagingSpec defines the desired state of Staging
type StagingSpec struct {
	// Optional human-readable description of this staging step
	Description string `json:"description,omitempty"`

	// Optional owner of the StagingSpec (distinct from the owner of the ingest or transform stages)
	Owner string `json:"owner,omitempty"`

	// Ingest defines how the pipeline is initiated, either by referencing a Kubernetes CronJob or Job or by referencing a Trigger CRD
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
// +kubebuilder:printcolumn:name="Description",type=string,JSONPath=`.spec.description`,description="Staging pipeline description"
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.owner`,description="Owner of the Staging",priority=1
// +kubebuilder:printcolumn:name="Ingest Kind",type=string,JSONPath=`.spec.ingest.kind`,description="Ingest type (cronjob, job, trigger)"
// +kubebuilder:printcolumn:name="Ingest Name",type=string,JSONPath=`.spec.ingest.name`,description="Name of the ingest resource"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Current Staging Status"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

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
