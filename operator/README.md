# Pipeline Forge Operator

A Kubernetes operator for orchestrating data pipeline stages with declarative configuration.

## Overview

The Pipeline Forge Operator is a Kubernetes-native controller that manages the lifecycle of data pipeline stages through Custom Resource Definitions (CRDs). It provides a declarative way to orchestrate data ingestion and transformation workflows within Kubernetes clusters.

## Architecture

The operator manages two primary custom resources:

- **Staging**: Represents a complete pipeline step that coordinates data ingestion and transformation
- **Trigger**: Provides event-driven capabilities for pipeline activation

### Staging Resource

A Staging resource represents a complete pipeline step that prepares data for analytics. It can:

- Reference existing Kubernetes CronJobs, Jobs, or Trigger resources
- Manage the lifecycle of ingestion containers
- Automatically trigger dbt transformations when ingestion completes
- Provide comprehensive observability and status tracking for pipeline stages
- Support retry mechanisms for both ingestion and transformation phases
- Suspend individual pipeline stages for maintenance or debugging
- Track detailed metrics including attempts, successes, and failures

### Trigger Resource

Trigger resources enable event-driven pipeline activation through:

- **GCS Triggers**: Monitor Google Cloud Storage buckets for file drops
- **Pub/Sub Triggers**: Listen to Google Cloud Pub/Sub messages with optional message filtering
- **BigQuery Triggers**: Watch for BigQuery table freshness and updates

Triggers support advanced features including:

- **Retry Mechanisms**: Configurable retry policies with exponential backoff
- **Cooldown Periods**: Prevent rapid re-triggering with configurable intervals
- **Run-Once Mode**: Execute triggers only once for one-time operations
- **Comprehensive Status Tracking**: Detailed monitoring of attempts, successes, and failures

## Installation

### Prerequisites

- Kubernetes cluster (v1.11.3+)
- kubectl configured to communicate with your cluster
- Docker or container runtime
- Go 1.24.0+ (for development)

### Deploy the Operator

1. **Build and push the operator image:**

```sh
make docker-build docker-push IMG=<your-registry>/pipeline-forge-operator:tag
```

2. **Install the CRDs:**

```sh
make install
```

3. **Deploy the operator:**

```sh
make deploy IMG=<your-registry>/pipeline-forge-operator:tag
```

4. **Apply sample resources:**

```sh
make apply-samples
```

### Uninstall

```sh
# Delete sample resources
make delete-samples

# Uninstall CRDs
make uninstall

# Remove operator deployment
make undeploy
```

## Examples

### Basic Staging with CronJob Reference

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: user-events-staging
  namespace: staging-sample
spec:
  description: "Loads user events from raw ingest into the analytics layer"
  owner: analytics-platform

  ingest:
    mode: reference # reference existing resource
    type: cronjob
    name: ingest-user-events # name of the resource to watch
    namespace: raw-jobs

  transform:
    name: user-events-transform
    project: data-group
    target: dev
    image: gcr.io/org/dbt-core:1.7.9
    models:
      - stg_user_events
    engine: dbt
    resources:
      requests:
        cpu: "250m"
        memory: "512Mi"
      limits:
        cpu: "1"
        memory: "2Gi"
```

### Event-Driven Staging with GCS Trigger

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Trigger
metadata:
  name: gcs-file-trigger
  namespace: staging-sample
spec:
  type: gcs
  name: gcs-file-trigger
  description: "Triggers pipeline when files are dropped in GCS bucket"
  owner: data-engineering

  gcs:
    bucket: "my-data-bucket"
    prefix: "exports/"
---
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: sales-data-staging
spec:
  ingest:
    mode: reference
    type: trigger
    name: gcs-file-trigger
    image: ghcr.io/org/sales-processor:latest

  transform:
    name: sales-data-transform
    project: analytics
    target: prod
    image: ghcr.io/org/dbt-core:latest
    models:
      - stg_sales_data
```

### BigQuery Data Freshness Monitoring

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Trigger
metadata:
  name: bq-freshness-trigger
  namespace: staging-sample
spec:
  type: bigquery
  name: bq-freshness-trigger
  description: "Monitors BigQuery table freshness and triggers pipeline updates"
  owner: data-engineering
  schedule: "*/15 * * * *" # Check every 15 minutes
  cooldownIntervalSeconds: 300 # 5 minute cooldown
  maxRetry: 3

  bigquery:
    project_id: "my-gcp-project"
    dataset_id: "analytics"
    table_id: "daily_sales"

---
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: sales-analytics-staging
spec:
  description: "Processes sales data when BigQuery table is updated"
  owner: analytics-team

  ingest:
    mode: reference
    type: trigger
    name: bq-freshness-trigger
    maxRetry: 2

  transform:
    name: sales-analytics-transform
    project: sales-analytics
    target: prod
    image: ghcr.io/org/dbt-core:latest
    models:
      - stg_daily_sales
      - fct_sales_metrics
    full_refresh: false
    maxRetry: 3
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"
```

## Advanced Features

### Retry Mechanisms

Both Staging and Trigger resources support configurable retry policies:

- **MaxRetry**: Set maximum retry attempts for failed operations
- **Exponential Backoff**: Automatic backoff between retry attempts
- **Status Tracking**: Monitor attempts, successes, and failures
- **Cooldown Periods**: Prevent rapid re-triggering with configurable intervals

### Status Monitoring

Comprehensive status tracking provides visibility into pipeline health:

- **Attempt Tracking**: Monitor total attempts, successful attempts, and failed attempts
- **Timing Information**: Track last attempt time, last failure time, and completion times
- **Job Monitoring**: Track Kubernetes Job names and execution status
- **Condition Reporting**: Kubernetes-style conditions for resource health

### Suspension and Control

Pipeline stages can be individually controlled:

- **Suspend Ingest**: Pause ingestion while keeping transformation active
- **Suspend Transform**: Pause transformation while ingestion continues
- **Run-Once Triggers**: Execute triggers only once for one-time operations

## Development

### Prerequisites

- Go 1.24.0+
- Docker
- kubectl
- Kind (for local testing)

### Setup Development Environment

1. **Clone and navigate to the operator directory:**

```sh
cd operator
```

2. **Install development dependencies:**

```sh
make manifests generate
```

3. **Run tests:**

```sh
make test
```

4. **Run e2e tests with Kind:**

```sh
make test-e2e
```

### Local Development

1. **Start a local Kind**:

```sh
make setup-test-e2e
```

2. **Install CRDs:**

```sh
make install
```

3. **Run the operator locally:**

```sh
make run
```

4. **Clean up local cluster**:

```sh
make cleanup-test-e2e
```

For more information on available make targets:

```sh
make help
```

## Metrics

Prometheus metrics are enabled for development and local validation. The operator
exposes a `pipeline_forge_reconciliations_total` counter partitioned by `result`:

| Result | Description |
|---|---|
| `initialized` | First reconciliation of a new Staging resource |
| `validation_failed` | Ingestion validation returned an error |
| `requeued` | Happy path, requeued after 60s |

To spin up kube-prometheus-stack locally and visualise metrics in Grafana:

```sh
make deploy-prometheus   # installs Prometheus + Grafana into the monitoring namespace

# Acess via port-forward
kubectl port-forward -n monitoring svc/kube-prometheus-kube-prome-prometheus 9090
kubectl port-forward -n monitoring svc/kube-prometheus-grafana 3000:80

make undeploy-prometheus # tears it down
```

Once running, query `pipeline_forge_reconciliations_total` in Grafana Explore.

## Related Components

- **Workloads**: See `../workloads/` for the actual data processing workloads
- **Main Project**: See `../README.md` for the complete Pipeline Forge project overview

## License

Copyright 2025 Daniel Blei.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
