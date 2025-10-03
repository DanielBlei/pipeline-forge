package v1alpha1

type StagingCondition string

const (
	StagingStatusCompleted     StagingCondition = "Completed"
	StagingConditionFailed     StagingCondition = "Failed"
	StagingConditionPending    StagingCondition = "Pending"
	StagingConditionRunning    StagingCondition = "Running"
	StagingConditionSuspended  StagingCondition = "Suspended"
	StagingConditionReady      StagingCondition = "Ready"
	StagingConditionInitiating StagingCondition = "Initiating"
)
