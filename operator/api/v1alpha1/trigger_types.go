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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements = corev1.ResourceRequirements

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
	//
	// +optional
	MessageFilter *PubSubMessageFilter `json:"messageFilter,omitempty"`
}

type PubSubMessageFilter struct {
	// Attribute name to match (e.g. "dataset")
	//
	// +kubebuilder:validation:Required
	Attribute string `json:"attribute"`

	// Required match value
	//
	// +kubebuilder:validation:Required
	Equals string `json:"equals"`
}

// BigQueryTriggerSpec defines a trigger based on BQ table freshness
type BigQueryTriggerSpec struct {
	// Google Cloud Project
	//
	// +kubebuilder:validation:Required
	Project string `json:"project_id"`

	// BigQuery Dataset ID
	//
	// +kubebuilder:validation:Required
	Dataset string `json:"dataset_id"`

	// BigQueryTable ID (e.g. "leads")
	//
	// +kubebuilder:validation:Required
	Table string `json:"table_id"`
}

// TriggerSpec defines the desired state of a Trigger resource.
type TriggerSpec struct {
	// Type specifies the kind of event that will activate this trigger.
	// Supported values: "gcs", "pubsub", "bigquery".
	// This field determines which trigger details must be provided.
	//
	// +kubebuilder:validation:Enum=gcs;pubsub;bigquery
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// Name of the trigger resource.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Description is an optional human-readable description of the trigger.
	//
	// +optional
	Description string `json:"description,omitempty"`

	// Owner is an optional field indicating the team or user responsible for this trigger.
	//
	// +optional
	Owner string `json:"owner,omitempty"`

	// Image is the container image used to run the trigger check.
	// This image is launched in a Kubernetes Job created by the Trigger controller.
	//
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Args are the arguments passed to the trigger container.
	// Useful for dynamic or parameterized behavior.
	//
	// +optional
	Args []string `json:"args,omitempty"`

	// Resources defines CPU and memory requests/limits for the trigger Job container.
	//
	// +optional
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Schedule is a cron-style schedule for when the trigger should check for new data or events.
	//
	// For example: "*/5 * * * *" (every 5 minutes)
	// Use standard cron syntax.
	//
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// CooldownIntervalSeconds defines the minimum time that must pass between successive runs, in seconds.
	// Prevents rapid re-triggering.
	// Example: 300 (for 5 minutes), 3600 (for 1 hour)
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	CooldownIntervalSeconds int32 `json:"cooldownIntervalSeconds,omitempty"`

	// RunOnce, if true, means this trigger should only run once.
	// Controller skips execution if status.lastTriggeredTime is set.
	//
	// +optional
	RunOnce bool `json:"runOnce,omitempty"`

	// GCS contains configuration for a GCS-based trigger.
	// Only set if Type is "gcs".
	//
	// +optional
	GCS *GCSTriggerSpec `json:"gcs,omitempty"`

	// PubSub contains configuration for a Pub/Sub-based trigger.
	// Only set if Type is "pubsub".
	//
	// +optional
	PubSub *PubSubTriggerSpec `json:"pubsub,omitempty"`

	// BigQuery contains configuration for a BigQuery-based trigger.
	// Only set if Type is "bigquery".
	//
	// +optional
	BigQuery *BigQueryTriggerSpec `json:"bigquery,omitempty"`

	// MaxRetry defines the maximum number of times to retry the transform step on failure.
	// If not set, defaults to 0 (no retries).
	//
	// +optional
	MaxRetry int32 `json:"maxRetry,omitempty"`

	// Suspend indicates whether the trigger is currently suspended.
	// If true, the controller will not process the trigger.
	//
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

// TriggerStatus defines the observed state of the Trigger resource.
// +kubebuilder:object:generate=true
type TriggerStatus struct {
	// Status is the status of the Trigger (e.g., Deployed, Failed, Pending, Running).
	Status *TriggerCondition `json:"status,omitempty"`

	// ObservedGeneration is the most recent generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represent the latest available observations of the trigger's state.
	// Examples: Ready, Failed, CooldownActive
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// Message provides human-readable information about the trigger status or last evaluation.
	Message string `json:"message,omitempty"`

	// LastTriggeredTime is the last time this trigger successfully resulted in a pipeline activation.
	LastTriggeredTime *metav1.Time `json:"lastTriggeredTime,omitempty"`

	// LastCheckTime is the last time the controller evaluated whether to trigger.
	LastCheckTime *metav1.Time `json:"lastCheckTime,omitempty"`

	// CooldownUntil is the next eligible time this trigger can fire, based on Cooldown setting.
	CooldownUntil *metav1.Time `json:"cooldownUntil,omitempty"`

	// LastRunJobName is the name of the last Kubernetes Job created by this trigger.
	LastRunJobName string `json:"lastRunJobName,omitempty"`

	// Attempts tracks the number of times this step has been attempted.
	Attempts int32 `json:"attempts,omitempty"`

	// SuccessfulAttempts tracks the number of successful attempts.
	SuccessfulAttempts int32 `json:"successfulAttempts,omitempty"`

	// FailedAttempts tracks the number of failed attempts.
	FailedAttempts int32 `json:"failedAttempts,omitempty"`

	// RetryCount tracks the number of retries for the current attempt.
	RetryCount int32 `json:"retryCount,omitempty"`

	// MaxRetries defines the maximum number of retries allowed.
	MaxRetries int32 `json:"maxRetries,omitempty"`

	// LastAttemptTime is the timestamp of the last attempt.
	LastAttemptTime *metav1.Time `json:"lastAttemptTime,omitempty"`

	// LastFailureTime is the timestamp of the last failure.
	LastFailureTime *metav1.Time `json:"lastFailureTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=trg
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`,description="Trigger type",priority=1
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.owner`,description="Owner",priority=1
// +kubebuilder:printcolumn:name="Last Triggered",type=date,JSONPath=`.status.lastTriggeredTime`,description="Last triggered time"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Status message"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Trigger is the Schema for the triggers API
type Trigger struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Trigger
	Spec TriggerSpec `json:"spec"`

	// status defines the observed state of Trigger
	Status TriggerStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// TriggerList contains a list of Trigger
type TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trigger `json:"items"`
}

// SetStatus is a helper function to set the trigger status condition
func (t *TriggerStatus) SetStatus(condition TriggerCondition) {
	t.Status = &condition
}

func init() {
	SchemeBuilder.Register(&Trigger{}, &TriggerList{})
}
