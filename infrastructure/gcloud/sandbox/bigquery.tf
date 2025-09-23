#-----------------------------------------------
# Ingest Workload (BigQuery)
#-----------------------------------------------

# Main ingest dataset for the pipeline-forge project (Raw Data)
resource "google_bigquery_dataset" "pipeline_forge_sandbox_ingest_dataset" {
  dataset_id    = "pipeline_forge_ingest"
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

# Transform dataset for the pipeline-forge project (DBT Models, e.g staging, mart, etc.)
resource "google_bigquery_dataset" "pipeline_forge_sandbox_transform_dataset" {
  dataset_id    = "pipeline_forge_transform"
  friendly_name = "Pipeline Forge Transform"
  description   = "Transform dataset for pipeline-forge data platform development"
  location      = "us-central1"

  delete_contents_on_destroy = true

  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "transform"
  }
}
