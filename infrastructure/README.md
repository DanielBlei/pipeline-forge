# Infrastructure

Cloud infrastructure automation and deployment tools for Pipeline Forge.

## Features

- **Terraform Modules** - Infrastructure as Code for cloud resources
- **Google Cloud Platform** - BigQuery, Secret Manager, and networking setup
- **Environment Management** - Dev, staging, and production configurations
- **Security Best Practices** - RBAC, network policies, and secret management

## Quick Start

```bash
cd gcloud/sandbox
terraform init
terraform plan
terraform apply
```

## What's Included

### **Google Cloud Platform Setup**

- **BigQuery Datasets** - Data warehouse for ingestion and transformation
- **Secret Manager** - Secure credential storage for database connections
- **IAM Configuration** - Service accounts and permissions
- **Networking** - VPC and firewall rules
- **Monitoring** - Cloud Logging and Monitoring setup

### **Terraform Modules**

- **BigQuery Module** - Dataset and table creation
- **IAM Module** - Service account and role management
- **Secrets Module** - Secret Manager integration
- **Networking Module** - VPC and firewall configuration

## Environment Configuration

### **Sandbox Environment**

```bash
cd gcloud/sandbox
terraform init
terraform plan -var-file="sandbox.tfvars"
terraform apply -var-file="sandbox.tfvars"
```

### **Production Environment**

```bash
cd gcloud/production
terraform init
terraform plan -var-file="production.tfvars"
terraform apply -var-file="production.tfvars"
```

## Infrastructure Components

### **BigQuery Setup**

- **Datasets**: `pipeline_forge_sandbox`, `pipeline_forge_prod`
- **Tables**: Configured for ingestion workloads
- **Access Control**: Service account-based permissions
- **Location**: Multi-region for high availability

### **Secret Manager**

- **Database Credentials** - MySQL and PostgreSQL passwords
- **API Keys** - External service integrations
- **Service Account Keys** - Application authentication
- **Rotation Policies** - Automated credential rotation

### **IAM Configuration**

- **Service Accounts**:
  - `pipeline-forge-ingest@project.iam.gserviceaccount.com`
  - `pipeline-forge-transform@project.iam.gserviceaccount.com`
  - `pipeline-forge-trigger@project.iam.gserviceaccount.com`
- **Roles**: BigQuery Admin, Secret Manager Accessor, Cloud Storage Admin

### **Networking**

- **VPC**: Isolated network for Pipeline Forge components
- **Firewall Rules**: Secure communication between services
- **Private Google Access**: Secure access to Google APIs
- **Cloud NAT**: Outbound internet access for workloads

## Security Best Practices

### **Access Control**

- **Principle of Least Privilege** - Minimal required permissions
- **Service Account Isolation** - Separate accounts for each workload
- **Resource-level Permissions** - Granular access control
- **Audit Logging** - Comprehensive access tracking

### **Secret Management**

- **Encryption at Rest** - All secrets encrypted with Google KMS
- **Rotation Policies** - Automated credential rotation
- **Access Logging** - Secret access monitoring
- **Version Control** - Secret versioning and rollback

### **Network Security**

- **Private Networks** - Isolated VPC for workloads
- **Firewall Rules** - Restrictive ingress/egress policies
- **Private Google Access** - Secure API communication
- **VPC Flow Logs** - Network traffic monitoring

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
# Create secrets
gcloud secrets create mysql-password --data-file=mysql-password.txt
gcloud secrets create postgres-password --data-file=postgres-password.txt

# Grant access to service accounts
gcloud secrets add-iam-policy-binding mysql-password \
  --member="serviceAccount:ingest@project.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### **3. Workload Deployment**

```bash
# Deploy Kubernetes operator
cd ../../operator
make deploy IMG=gcr.io/project/pipeline-forge-operator:latest

# Deploy workloads
kubectl apply -f config/samples/
```

## Monitoring and Observability

### **Cloud Logging**

- **Structured Logging** - JSON-formatted logs from all workloads
- **Log Aggregation** - Centralized log collection
- **Log Retention** - Configurable retention policies
- **Log Analysis** - Cloud Logging queries and dashboards

### **Cloud Monitoring**

- **Custom Metrics** - Pipeline execution metrics
- **Alerting** - Automated alerts for failures
- **Dashboards** - Real-time pipeline monitoring
- **Uptime Checks** - Service availability monitoring

## Cost Optimization

### **Resource Management**

- **Auto-scaling** - Dynamic resource allocation
- **Resource Quotas** - Cost control and limits
- **Scheduling** - Non-production environment shutdown
- **Monitoring** - Cost tracking and optimization

### **BigQuery Optimization**

- **Partitioning** - Table partitioning for cost reduction
- **Clustering** - Query performance optimization
- **Slots Management** - Flexible slot allocation
- **Data Lifecycle** - Automated data retention policies

## Troubleshooting

### **Common Issues**

1. **Permission Denied**

   ```bash
   # Check service account permissions
   gcloud projects get-iam-policy PROJECT_ID
   ```

2. **Secret Access Issues**

   ```bash
   # Verify secret access
   gcloud secrets versions access latest --secret="SECRET_NAME"
   ```

3. **BigQuery Connection Issues**
   ```bash
   # Test BigQuery access
   bq ls --project_id=PROJECT_ID
   ```

### **Debug Commands**

```bash
# Check Terraform state
terraform show

# Validate configuration
terraform validate

# Check resource status
gcloud compute instances list
gcloud sql instances list
```

## Environment Variables

| Variable            | Description            | Example                  |
| ------------------- | ---------------------- | ------------------------ |
| `GOOGLE_PROJECT_ID` | GCP Project ID         | `pipeline-forge-sandbox` |
| `GOOGLE_REGION`     | GCP Region             | `us-central1`            |
| `GOOGLE_ZONE`       | GCP Zone               | `us-central1-a`          |
| `BIGQUERY_DATASET`  | BigQuery Dataset       | `pipeline_forge_sandbox` |
| `SECRET_PROJECT_ID` | Secret Manager Project | `pipeline-forge-secrets` |
