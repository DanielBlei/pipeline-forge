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

type IngestMode string

const (
	IngestModeReference IngestMode = "reference"
	IngestModeBootstrap IngestMode = "bootstrap"
)

// IngestSpec defines how this staging process is triggered by an upstream data ingestion step.
//
// The controller does not run the ingestion itself. Instead, it watches for the completion
// of an external Kubernetes Job, CronJob, or a Trigger CRD.
//
// - If Mode is "reference", the referenced resource must already exist.
// - If Mode is "bootstrap", the controller may create and manage the ingestion resource (future support).
type IngestSpec struct {
	// Mode determines whether the ingestion source is a reference to an existing resource,
	// or should be created and managed by this Staging resource in the future.
	//
	// Accepted values:
	// - "reference": reference an existing CronJob, Job, or Trigger
	// - "bootstrap": placeholder for future functionality to create the ingestion job
	//
	// +kubebuilder:validation:Enum=reference;bootstrap
	// +kubebuilder:validation:Required
	Mode IngestMode `json:"mode"`

	// Type specifies the kind of resource used to signal ingestion readiness.
	//
	// Accepted values:
	// - "cronjob": watches for successful Job completions from a CronJob
	// - "job": watches a single-use Kubernetes Job
	// - "trigger": watches a custom Trigger resource (e.g., based on PubSub, GCS, BigQuery)
	//
	// +kubebuilder:validation:Enum=cronjob;job;trigger
	// +kubebuilder:validation:Required
	Type IngestType `json:"type"`

	// Name of the ingestion resource.
	// This must match the Kubernetes or custom resource to be watched.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace where the ingestion resource is defined.
	// If omitted, defaults to the namespace of the Staging object.
	//
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Image is the container image used to run the ingest. (only relevant when mode is "bootstrap")
	// This image is launched in a Kubernetes from a Job or Cronjob created by the Staging controller.
	//
	// +optional
	Image string `json:"image"`

	// Args specifies the command-line arguments to pass to the container,
	// overriding the default container command (entrypoint).
	// If set, these arguments will replace the container's default command.
	//
	// +optional
	Args []string `json:"args,omitempty"`

	// Resources defines CPU and memory requests/limits for the Cronjob/Job container.
	// If Type is "trigger", this field is ignored.
	//
	// +optional
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Schedule defines the cron schedule (in standard cron format) for the ingestion job.
	// This field is only relevant when Mode is "bootstrap" and Type is "cronjob".
	// If Mode is "reference" or Type is not "cronjob", this field is ignored.
	//
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// PollIntervalSeconds specifies how frequently (in seconds) the controller checks for changes
	// or status updates in the ingestion resource.
	// If not set, defaults to 120 seconds.
	//
	// +optional
	PollIntervalSeconds int32 `json:"pollIntervalSeconds,omitempty"`

	// Suspend indicates whether the ingestion action is currently suspended.
	// If true, the controller will not process or watch the referenced ingestion resource.
	//
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// Owner identifies the team or user responsible for this ingestion step.
	//
	// +optional
	Owner string `json:"owner,omitempty"`

	// MaxRetry defines the maximum number of times to retry the transform step on failure.
	// If not set, defaults to 0 (no retries).
	//
	// +optional
	MaxRetry int32 `json:"maxRetry,omitempty"`
}

// TransformSpec defines the transformation logic to run after ingestion.
//
// This stage creates a Kubernetes Job that runs the specified container image
// (e.g., a dbt-core image) to execute the transformation. While this is initially
// designed for dbt workflows, the spec is extensible to other engines such as
// Spark, Python, or custom tools in the future.
type TransformSpec struct {
	// Name is a reference or descriptive label for this transform step.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Project is the name of the transformation project (e.g., dbt project name).
	//
	// +kubebuilder:validation:Required
	Project string `json:"project"`

	// Target defines the execution target or environment profile (e.g., "dev", "prod").
	//
	// +kubebuilder:validation:Required
	Target string `json:"target"`

	// Image is the container image used to run the transform.
	// This image is launched in a Kubernetes Job created by the Staging controller.
	//
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Args specifies the command-line arguments to pass to the container,
	// overriding the default container command (entrypoint).
	// If set, these arguments will replace the container's default command.
	//
	// +optional
	Args []string `json:"args,omitempty"`

	// Models defines the transformation model names to run.
	// For dbt, these are model names. For other engines, these could be script names, tasks, etc.
	//
	// +optional
	Models []string `json:"models,omitempty"`

	// Resources defines CPU and memory requests/limits for the Job container.
	//
	// +optional
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Engine defines the transformation engine to use (e.g., "dbt", "spark", "custom").
	// Defaults to "dbt" if not specified.
	//
	// +kubebuilder:validation:Enum=dbt;spark;custom
	// +kubebuilder:default=dbt
	Engine TransformEngineType `json:"engine,omitempty"`

	// Owner identifies the team or individual responsible for the transformation.
	//
	// +optional
	Owner string `json:"owner,omitempty"`

	// Suspend indicates whether the Transform action is currently suspended.
	// If true, the controller will not process the transform step(s)
	//
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// Whether to force a full refresh of data (e.g., `dbt run --full-refresh`).
	// Only applicable to engines that support it (like dbt).
	//
	// +optional
	FullRefresh bool `json:"full_refresh,omitempty"`

	// MaxRetry defines the maximum number of times to retry the transform step on failure.
	// If not set, defaults to 0 (no retries).
	//
	// +optional
	MaxRetry int32 `json:"maxRetry,omitempty"`
}

// StagingSpec defines the desired state of Staging
type StagingSpec struct {
	// Optional human-readable description of this staging step
	Description string `json:"description,omitempty"`

	// Optional owner of the StagingSpec (distinct from the owner of the ingest or transform stages)
	Owner string `json:"owner,omitempty"`

	// Ingest defines how the pipeline is initiated, either by referencing a Kubernetes CronJob or Job or by referencing a Trigger CRD
	//
	// +kubebuilder:validation:Required
	Ingest IngestSpec `json:"ingest"`

	// Transform defines the dbt-based transformation logic
	//
	// +kubebuilder:validation:Required
	Transform TransformSpec `json:"transform"`
}

// InternalStatus represents the current state of either the ingestion or transform step in a pipeline.
// +kubebuilder:object:generate=true
type InternalStatus struct {
	// Status of the associated staging object (e.g., Complete, Running, Failed, Pending, Unknown).
	Status *StagingCondition `json:"status,omitempty"`

	// Message provides additional information or error messages about the step status.
	Message string `json:"message,omitempty"`

	// LastRunJobName is the name of the last Kubernetes Job created by this resource.
	LastRunJobName string `json:"lastRunJobName,omitempty"`

	// LastCompletedTime is the last time the step (e.g., ingestion CronJob) completed.
	LastCompletedTime metav1.Time `json:"lastCompletedTime,omitempty"`

	// LastCheckedTime is the last time the status was checked.
	LastCheckedTime metav1.Time `json:"lastCheckedTime,omitempty"`

	// LastAttemptTime is the timestamp of the last attempt.
	LastAttemptTime *metav1.Time `json:"lastAttemptTime,omitempty"`

	// LastFailureTime is the timestamp of the last failure.
	LastFailureTime *metav1.Time `json:"lastFailureTime,omitempty"`

	// Attempts tracks the number of times this step has been attempted.
	Attempts int32 `json:"attempts,omitempty"`

	// RetryCount tracks the number of retries for the current attempt.
	RetryCount int32 `json:"retryCount,omitempty"`

	// MaxRetries defines the maximum number of retries allowed.
	MaxRetries int32 `json:"maxRetries,omitempty"`

	// SuccessfulAttempts tracks the number of successful attempts.
	SuccessfulAttempts int32 `json:"successfulAttempts,omitempty"`

	// FailedAttempts tracks the number of failed attempts.
	FailedAttempts int32 `json:"failedAttempts,omitempty"`
}

// StagingStatus defines the observed state of a Staging resource.
type StagingStatus struct {
	// Status is the status of the Staging (e.g., Deployed, Failed, Pending, Running).
	Status *StagingCondition `json:"status,omitempty"`

	// ObservedGeneration is the most recent generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Ingest is the status of the ingestion step.
	Ingest InternalStatus `json:"ingest,omitempty"`

	// Transform is the status of the transformation step.
	Transform InternalStatus `json:"transform,omitempty"`

	// Conditions is a list of status conditions for the Staging resource.
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=stg
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Overall staging status",priority=1
// +kubebuilder:printcolumn:name="Owner",type=string,JSONPath=`.spec.owner`,description="Owner",priority=1
// +kubebuilder:printcolumn:name="Ingest",type=string,JSONPath=`.spec.ingest.type`,description="Ingest type"
// +kubebuilder:printcolumn:name="Ingest Status",type=string,JSONPath=`.status.ingest.status`,description="Ingest status"
// +kubebuilder:printcolumn:name="Transform Status",type=string,JSONPath=`.status.transform.status`,description="Transform status"
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

// SetStatus is a helper function to set the staging status condition
func (s *StagingStatus) SetStatus(condition StagingCondition) {
	s.Status = &condition
}

func init() {
	SchemeBuilder.Register(&Staging{}, &StagingList{})
}
