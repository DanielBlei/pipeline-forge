terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }

  # Keep track of the state in the GCS bucket
  backend "gcs" {
    bucket = "pipeline-forge-terraform-state"
    prefix = "sandbox"
  }
}

provider "google" {
  project = "pipeline-forge"
  region  = "us-central1"
}
