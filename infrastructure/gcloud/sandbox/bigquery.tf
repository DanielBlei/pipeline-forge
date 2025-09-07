#-----------------------------------------------
# Ingest Workload (BigQuery)
#-----------------------------------------------

# Main dataset for the pipeline-forge project
resource "google_bigquery_dataset" "pipeline_forge_sandbox_dataset" {
  dataset_id    = "pipeline_forge_sandbox"
  friendly_name = "Pipeline Forge"
  description   = "Sandbox dataset for pipeline-forge data platform development"
  location      = "us-central1"

  delete_contents_on_destroy = true

  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "ingest"
  }
}
