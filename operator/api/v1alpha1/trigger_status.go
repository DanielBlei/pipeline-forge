package v1alpha1

type TriggerCondition string

const (
	TriggerStatusCompleted     TriggerCondition = "Completed"
	TriggerConditionFailed     TriggerCondition = "Failed"
	TriggerConditionPending    TriggerCondition = "Pending"
	TriggerConditionRunning    TriggerCondition = "Running"
	TriggerConditionSuspended  TriggerCondition = "Suspended"
	TriggerConditionReady      TriggerCondition = "Ready"
	TriggerConditionInitiating TriggerCondition = "Initiating"
)
