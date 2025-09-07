#-----------------------------------------------
# Ingest Workload (IAM)
#-----------------------------------------------

# Service Account
resource "google_service_account" "ingest_workload" {
  account_id   = "ingest-workload"
  display_name = "Ingest Workload"
  description  = "Service Account for Ingest Workload"
}

# Grant Dataset Owner to Service Account
resource "google_bigquery_dataset_iam_member" "pipeline_forge_sandbox_dataset_owner" {
  dataset_id = google_bigquery_dataset.pipeline_forge_sandbox_dataset.dataset_id
  role       = "roles/bigquery.dataOwner"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}

# Grant Secret Manager Accessor to Service Account
resource "google_secret_manager_secret_iam_member" "postgres_secret_manager_password_accessor" {
  secret_id  = google_secret_manager_secret.postgres_password.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}

# Grant Secret Manager Accessor to Service Account
resource "google_secret_manager_secret_iam_member" "mysql_secret_manager_password_accessor" {
  secret_id  = google_secret_manager_secret.mysql_password.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}
