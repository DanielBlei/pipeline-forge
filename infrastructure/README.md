# Infrastructure

Cloud infrastructure automation and deployment tools for Pipeline Forge.

## Features

- **Terraform Infrastructure** - Infrastructure as Code for cloud resources
- **Google Cloud Platform** - BigQuery, Secret Manager, Pub/Sub, and Cloud Storage
- **Sandbox Environment** - Development environment for Pipeline Forge workloads
- **Security Best Practices** - Service account isolation and secret management

## Quick Start

```bash
cd gcloud/sandbox
terraform init
terraform plan
terraform apply
```

## What's Included

### **Google Cloud Platform Setup**

- **BigQuery Datasets** - `pipeline_forge_ingest` and `pipeline_forge_transform` datasets
- **Secret Manager** - Secure credential storage for MySQL and PostgreSQL connections
- **IAM Configuration** - Dedicated service accounts for each workload
- **Pub/Sub** - Topic and subscription for trigger workload messaging
- **Cloud Storage** - GCS bucket for trigger workload data

### **Terraform Resources**

- **BigQuery Resources** - Dataset creation and IAM bindings
- **IAM Resources** - Service account and role management
- **Secret Manager Resources** - Secret creation and access control
- **Pub/Sub Resources** - Topic and subscription configuration
- **Storage Resources** - GCS bucket creation and permissions

## Environment Configuration

### **Sandbox Environment**

```bash
cd gcloud/sandbox
terraform init
terraform plan
terraform apply
```

## Infrastructure Components

### **BigQuery Setup**

- **Datasets**: `pipeline_forge_ingest`, `pipeline_forge_transform`
- **Access Control**: Service account-based permissions
- **Location**: `us-central1` region

### **Secret Manager**

- **Database Credentials** - MySQL and PostgreSQL passwords for sandbox environment
- **Access Control** - Service account-based secret access

### **IAM Configuration**

- **Service Accounts**:
  - `ingest-workload@pipeline-forge.iam.gserviceaccount.com`
  - `transform-workload@pipeline-forge.iam.gserviceaccount.com`
  - `trigger-workload@pipeline-forge.iam.gserviceaccount.com`
- **Roles**: BigQuery Data Owner, Secret Manager Accessor, Storage Object Owner, Pub/Sub Subscriber

### **Pub/Sub**

- **Topic**: `pipeline_forge_topic` for trigger workload messaging
- **Subscription**: `pipeline_forge_subscription` with acknowledge-based message handling
- **Access Control**: Service account-based subscriber permissions

### **Cloud Storage**

- **Bucket**: `pipeline_forge_bucket` for trigger workload data storage
- **Location**: `us-central1` region with standard storage class
- **Access Control**: Service account-based object ownership

## Security Best Practices

### **Access Control**

- **Principle of Least Privilege** - Minimal required permissions
- **Service Account Isolation** - Separate accounts for each workload
- **Resource-level Permissions** - Granular access control

### **Secret Management**

- **Encryption at Rest** - All secrets encrypted with Google KMS
- **Access Control** - Service account-based secret access

## Deployment Workflow

### **1. Infrastructure Provisioning**

```bash
# Initialize Terraform
terraform init

# Plan infrastructure changes
terraform plan

# Apply changes
terraform apply
```

### **2. Secret Configuration**

```bash
# Create secrets (handled by Terraform)
# MySQL and PostgreSQL passwords are created automatically
# Service account access is configured via Terraform IAM bindings
```

### **3. Workload Deployment**

```bash
# Deploy Kubernetes operator (if needed)
cd ../../operator
make deploy IMG=gcr.io/project/pipeline-forge-operator:latest

# Deploy workloads
kubectl apply -f config/samples/
```

## Monitoring and Observability

### **Cloud Logging**

- **Structured Logging** - JSON-formatted logs from all workloads
- **Log Aggregation** - Centralized log collection
- **Log Analysis** - Cloud Logging queries and dashboards

### **Cloud Monitoring**

- **Custom Metrics** - Pipeline execution metrics
- **Alerting** - Automated alerts for failures
- **Dashboards** - Real-time pipeline monitoring

## Troubleshooting

### **Common Issues**

1. **Permission Denied**

   ```bash
   # Make sure to login with the correct account (e.g. Default account)
   gcloud auth application-default login

   # Or Impersonate a service account (e.g. ingest-workload)
   gcloud auth application-default login --impersonate-service-account=ingest-workload@pipeline-forge.iam.gserviceaccount.com
   ```

2. **Secret Access Issues**

   ```bash
   # Verify secret access
   gcloud secrets versions access latest --secret="SECRET_NAME"
   ```

3. **BigQuery Connection Issues**
   ```bash
   # Test BigQuery access
   bq ls --project_id=pipeline-forge
   ```

## Environment Variables

| Variable                     | Description       | Example                    |
| ---------------------------- | ----------------- | -------------------------- |
| `GOOGLE_PROJECT_ID`          | GCP Project ID    | `pipeline-forge`           |
| `GOOGLE_REGION`              | GCP Region        | `us-central1`              |
| `BIGQUERY_INGEST_DATASET`    | Ingest Dataset    | `pipeline_forge_ingest`    |
| `BIGQUERY_TRANSFORM_DATASET` | Transform Dataset | `pipeline_forge_transform` |
| `PUBSUB_TOPIC`               | Pub/Sub Topic     | `pipeline_forge_topic`     |
| `GCS_BUCKET`                 | Storage Bucket    | `pipeline_forge_bucket`    |
