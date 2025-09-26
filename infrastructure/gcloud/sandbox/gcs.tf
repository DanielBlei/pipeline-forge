#-----------------------------------------------
# Trigger Workload (GCS)
#-----------------------------------------------

# GCS bucket for the trigger workload
resource "google_storage_bucket" "trigger_workload_bucket" {
  name          = "pipeline_forge_bucket"
  location      = "us-central1"
  storage_class = "STANDARD"
  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "trigger"
  }
}
