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
- Provide observability and status tracking for pipeline stages

### Trigger Resource

Trigger resources enable event-driven pipeline activation through:

- **GCS Triggers**: Monitor Google Cloud Storage buckets for file drops
- **Pub/Sub Triggers**: Listen to Google Cloud Pub/Sub messages
- **BigQuery Triggers**: Watch for BigQuery table updates

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
    image: ghcr.io/org/ingest:latest
    resources:
      requests:
        cpu: "100m"
        memory: "256Mi"
      limits:
        cpu: "500m"
        memory: "1Gi"

  transform:
    name: user-events-transform
    project: data-group
    target: dev
    image: ghcr.io/org/dbt-core:1.7.9
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
