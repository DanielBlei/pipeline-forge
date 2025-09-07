#-----------------------------------------------
# Ingest Workload (IAM)
#-----------------------------------------------

# Database secrets for local development
resource "google_secret_manager_secret" "postgres_password" {
  secret_id = "docker-compose-postgres-password"

  replication {
    auto {}
  }

  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "ingest"
  }
}

resource "google_secret_manager_secret" "mysql_password" {
  secret_id = "docker-compose-mysql-password"
  replication {
    auto {}
  }
  labels = {
    project     = "pipeline-forge"
    environment = "sandbox"
    workload    = "ingest"
  }
}

# Secret versions - reads from .env file with fallback to variables
resource "google_secret_manager_secret_version" "postgres_password" {
  secret = google_secret_manager_secret.postgres_password.id
  secret_data = try(
    trimspace(regex("POSTGRES_PASSWORD=(.+)", file("../../../dev/.env"))),
    var.postgres_password
  )
  depends_on = [google_secret_manager_secret.postgres_password]
}

resource "google_secret_manager_secret_version" "mysql_password" {
  secret = google_secret_manager_secret.mysql_password.id
  secret_data = try(
    trimspace(regex("MYSQL_ROOT_PASSWORD=(.+)", file("../../../dev/.env"))),
    var.mysql_password
  )
  depends_on = [google_secret_manager_secret.mysql_password]
}
