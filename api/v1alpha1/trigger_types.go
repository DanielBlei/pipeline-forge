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

// GCSTriggerSpec defines a trigger based on GCS file drops
type GCSTriggerSpec struct {
	// GCS bucket name (without gs:// prefix)
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`

	// Optional prefix to match files (e.g. "exports/")
	Prefix string `json:"prefix,omitempty"`
}

// PubSubTriggerSpec defines a trigger from a Pub/Sub message
type PubSubTriggerSpec struct {
	// Pub/Sub topic to subscribe to
	// +kubebuilder:validation:Required
	Topic string `json:"topic"`

	// Optional message filter
	// +optional
	MessageFilter *PubSubMessageFilter `json:"messageFilter,omitempty"`
}

type PubSubMessageFilter struct {
	// Attribute name to match (e.g. "dataset")
	// +kubebuilder:validation:Required
	Attribute string `json:"attribute"`

	// Required match value
	// +kubebuilder:validation:Required
	Equals string `json:"equals"`
}

// BigQueryTriggerSpec defines a trigger based on BQ table freshness
type BigQueryTriggerSpec struct {
	// Google Cloud Project
	// +kubebuilder:validation:Required
	Project string `json:"project_id"`

	// BigQuery Dataset ID
	// +kubebuilder:validation:Required
	Dataset string `json:"dataset_id"`

	// BigQueryTable ID (e.g. "leads")
	// +kubebuilder:validation:Required
	Table string `json:"table_id"`
}

// TriggerSpec defines the mechanism for activating a pipeline stage,
// either on a schedule or in response to an event (GCS, Pub/Sub, BQ, etc.).
type TriggerSpec struct {
	// Type determines the trigger mechanism: cron, gcs, pubsub, bigquery, etc.
	// +kubebuilder:validation:Enum=gcs;pubsub;bigquery
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// +optional
	GCS *GCSTriggerSpec `json:"gcs,omitempty"`

	// +optional
	PubSub *PubSubTriggerSpec `json:"pubsub,omitempty"`

	// +optional
	BigQuery *BigQueryTriggerSpec `json:"bigquery,omitempty"`
}

// TriggerStatus defines the observed state of Trigger.
type TriggerStatus struct {
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Kind",type=string,JSONPath=`.spec.kind`,description="Type of trigger (gcs, pubsub, bigquery)"
// +kubebuilder:printcolumn:name="GCS Bucket",type=string,JSONPath=`.spec.gcs.bucket`,description="GCS bucket (if GCS trigger)",priority=1
// +kubebuilder:printcolumn:name="PubSub Topic",type=string,JSONPath=`.spec.pubsub.topic`,description="PubSub topic (if PubSub trigger)",priority=1
// +kubebuilder:printcolumn:name="BQ Dataset",type=string,JSONPath=`.spec.bigquery.dataset`,description="BigQuery dataset (if BQ trigger)",priority=1
// +kubebuilder:printcolumn:name="BQ Table",type=string,JSONPath=`.spec.bigquery.table`,description="BigQuery table (if BQ trigger)",priority=1
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Current status"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Trigger is the Schema for the triggers API
type Trigger struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Trigger
	// +required
	Spec TriggerSpec `json:"spec"`

	// status defines the observed state of Trigger
	// +optional
	Status TriggerStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// TriggerList contains a list of Trigger
type TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trigger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Trigger{}, &TriggerList{})
}
