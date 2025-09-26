#-----------------------------------------------
# Trigger Workload (PubSub)
#-----------------------------------------------

# PubSub topic for the trigger workload
resource "google_pubsub_topic" "pipeline_forge_topic" {
  name = "pipeline_forge_topic"
  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "trigger"
  }
}

# PubSub subscription with "acknowledge" (ack) deletes messages after they are read
resource "google_pubsub_subscription" "pipeline_forge_subscription" {
  name  = "pipeline_forge_subscription"
  topic = google_pubsub_topic.pipeline_forge_topic.id
  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "trigger"
  }
}
