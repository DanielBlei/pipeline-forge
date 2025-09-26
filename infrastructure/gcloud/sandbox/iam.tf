#-----------------------------------------------
# Ingest Workload (IAM)
#-----------------------------------------------

# Service Account
resource "google_service_account" "ingest_workload" {
  account_id   = "ingest-workload"
  display_name = "Ingest Workload"
  description  = "Service Account for Ingest Workload"
}

# Grant Dataset Owner to Service Account for the Ingest Dataset
resource "google_bigquery_dataset_iam_member" "pipeline_forge_sandbox_dataset_owner_ingest" {
  dataset_id = google_bigquery_dataset.pipeline_forge_sandbox_ingest_dataset.dataset_id
  role       = "roles/bigquery.dataOwner"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}

# Grant Secret Manager Accessor for Postgres Secret to Service Account
resource "google_secret_manager_secret_iam_member" "postgres_secret_manager_password_accessor" {
  secret_id  = google_secret_manager_secret.postgres_password.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}

# Grant Secret Manager Accessor for MYSQL Secret to Service Account
resource "google_secret_manager_secret_iam_member" "mysql_secret_manager_password_accessor" {
  secret_id  = google_secret_manager_secret.mysql_password.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.ingest_workload.email}"
  depends_on = [google_service_account.ingest_workload]
}

# -----------------------------------------------
# Transform Workload (IAM)
# -----------------------------------------------

# Service Account
resource "google_service_account" "transform_workload" {
  account_id   = "transform-workload"
  display_name = "Transform Workload"
  description  = "Service Account for Transform Workload"
}

# Grant Dataset Owner to Service Account for the Transform Dataset
resource "google_bigquery_dataset_iam_member" "pipeline_forge_sandbox_dataset_owner_transform" {
  dataset_id = google_bigquery_dataset.pipeline_forge_sandbox_transform_dataset.dataset_id
  role       = "roles/bigquery.dataOwner"
  member     = "serviceAccount:${google_service_account.transform_workload.email}"
  depends_on = [google_service_account.transform_workload]
}

# -----------------------------------------------
# Trigger Workload (IAM)
# -----------------------------------------------

# Service Account
resource "google_service_account" "pipeline_forge_workload_sa" {
  account_id   = "trigger-workload"
  display_name = "Trigger Workload"
  description  = "Service Account for Trigger Workload"
}

# IAM binding to allow a service account to pull (read) messages from the subscription
resource "google_pubsub_subscription_iam_member" "pipeline_forge_reader" {
  subscription = google_pubsub_subscription.pipeline_forge_subscription.id
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${google_service_account.pipeline_forge_workload_sa.email}"
  depends_on   = [google_service_account.pipeline_forge_workload_sa]
}

# Grant Storage Object Viewer to Service Account
resource "google_storage_bucket_iam_member" "pipeline_forge_bucket_owner" {
  bucket     = google_storage_bucket.pipeline_forge_bucket.name
  role       = "roles/storage.objectOwner"
  member     = "serviceAccount:${google_service_account.pipeline_forge_workload_sa.email}"
  depends_on = [google_service_account.pipeline_forge_workload_sa]
}
