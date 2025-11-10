package v1alpha1

// StatusType represents the status of a resource
type StatusType string

// Status constants for consistent status values across all resources
const (
	StatusPending   StatusType = "Pending"
	StatusRunning   StatusType = "Running"
	StatusCompleted StatusType = "Completed"
	StatusFailed    StatusType = "Failed"
	StatusSuspended StatusType = "Suspended"
	StatusUnknown   StatusType = "Unknown"
	StatusReady     StatusType = "Ready"
)

// EngineType represents the transformation engine
type TransformEngineType string

const (
	EngineDBT    TransformEngineType = "dbt"
	EngineSpark  TransformEngineType = "spark"
	EngineCustom TransformEngineType = "custom"
)

// TriggerType represents the type of trigger
type TriggerType string

const (
	TriggerTypeGCS      TriggerType = "gcs"
	TriggerTypePubSub   TriggerType = "pubsub"
	TriggerTypeBigQuery TriggerType = "bigquery"
)

// IngestType represents the type of ingestion
type IngestType string

const (
	IngestTypeCronjob IngestType = "cronjob"
	IngestTypeJob     IngestType = "job"
	IngestTypeTrigger IngestType = "trigger"
)

// ObjectCondition represents the condition of an object (Ingest, Transform, Trigger, etc.)
type ObjCondition string

const (
	ObjectStatusCompleted  ObjCondition = "Completed"
	ObjConditionFailed     ObjCondition = "Failed"
	ObjConditionPending    ObjCondition = "Pending"
	ObjConditionRunning    ObjCondition = "Running"
	ObjConditionSuspended  ObjCondition = "Suspended"
	ObjConditionReady      ObjCondition = "Ready"
	ObjConditionInitiating ObjCondition = "Initiating"
)
