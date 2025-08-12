# Pipeline Forge 🔥

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.11+-blue.svg)](https://kubernetes.io/)
[![Python](https://img.shields.io/badge/Python-3.13+-blue.svg)](https://python.org/)
[![dbt](https://img.shields.io/badge/dbt-Core-blue.svg)](https://www.getdbt.com/)

> A Kubernetes-native platform for building modern, declarative data pipelines with clear boundaries between ingestion and transformation.

## 🚀 Overview

Pipeline Forge is a complete solution for orchestrating data pipelines in Kubernetes environments. It combines a powerful Kubernetes operator with specialized workloads to provide a declarative, event-driven approach to data pipeline management.

### 🎯 The Problem

Modern data teams face challenges with:

- **Complex Orchestration**: Managing dependencies between data ingestion and transformation
- **Event-Driven Requirements**: Responding to file drops, database changes, and streaming events
- **Infrastructure Complexity**: Deploying and scaling data processing workloads
- **Observability Gaps**: Tracking pipeline health and data lineage
- **Team Coordination**: Coordinating between data engineering and platform teams

### 💡 The Solution

Pipeline Forge provides a Kubernetes-native platform that:

- **Declaratively Orchestrates** data pipelines using Custom Resource Definitions (CRDs)
- **Event-Driven Architecture** responds to GCS file drops, Pub/Sub messages, and BigQuery updates
- **Clear Separation** between data ingestion and transformation phases
- **Built-in Observability** with comprehensive status tracking and monitoring
- **GitOps Ready** configuration that fits modern deployment practices

## 🏗️ Architecture

Pipeline Forge consists of two main components that work together seamlessly:

### 🎛️ Kubernetes Operator

The operator manages the lifecycle of data pipeline stages through:

- **Staging Resources**: Complete pipeline steps that coordinate ingestion and transformation
- **Trigger Resources**: Event-driven activation for pipelines
- **Automatic Reconciliation**: Ensures pipeline state matches desired configuration

### 🔧 Specialized Workloads

Pre-built, production-ready data processing components:

- **Ingest Workloads**: Type-safe data ingestion from MySQL, PostgreSQL, and more
- **Transform Workloads**: dbt-based data transformation with version control
- **Trigger Workloads**: Event processing for GCS, Pub/Sub, and BigQuery

## 🛠️ Technology Stack

| Component          | Technology                    | Purpose                                   |
| ------------------ | ----------------------------- | ----------------------------------------- |
| **Operator**       | Go, Kubernetes, Kubebuilder   | Pipeline orchestration and CRD management |
| **Ingest**         | Python 3.13+, Pydantic, Typer | Type-safe data ingestion with validation  |
| **Transform**      | dbt Core, BigQuery            | Data transformation and analytics         |
| **Triggers**       | GCS, Pub/Sub, BigQuery APIs   | Event-driven pipeline activation          |
| **Infrastructure** | Kubernetes, Docker, Helm      | Container orchestration and deployment    |

## ✨ Key Features

### 📝 Declarative Pipeline Configuration

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: user-events-pipeline
spec:
  ingest:
    mode: reference
    type: cronjob
    name: user-events-ingest
    image: ghcr.io/org/ingest:latest

  transform:
    name: user-events-transform
    project: analytics
    target: prod
    image: ghcr.io/org/dbt-core:latest
    models:
      - stg_user_events
      - fct_user_events
```

### ⚡ Event-Driven Orchestration

- **GCS Triggers**: Monitor bucket changes and trigger pipelines
- **Pub/Sub Triggers**: React to real-time messages
- **BigQuery Triggers**: Watch for table updates and freshness
- **Custom Events**: Extensible trigger system for any event source

### 🔄 Clear Pipeline Boundaries

- **Ingestion Phase**: Data extraction and loading with specialized workloads
- **Transformation Phase**: dbt-based data modeling and analytics
- **Observability**: Built-in status tracking and health monitoring

### ☸️ Kubernetes-Native

- **CRD-Based**: Native Kubernetes resources for pipeline definition
- **RBAC Integration**: Fine-grained access control for teams
- **Resource Management**: CPU, memory, and storage allocation
- **Scaling**: Automatic scaling based on workload demands

## 🔄 Workflow

### 1. 📋 Pipeline Definition

Define your data pipeline using declarative YAML configuration:

- Specify data sources and ingestion schedules
- Define transformation models and dependencies
- Configure event triggers for real-time processing

### 2. 🎛️ Operator Orchestration

The operator automatically:

- Creates and manages Kubernetes Jobs and CronJobs
- Monitors pipeline execution and status
- Handles failures and retries
- Triggers transformations when ingestion completes

### 3. ⚙️ Workload Execution

Specialized workloads handle the actual data processing:

- **Ingest**: Extract data from sources with type-safe configuration
- **Transform**: Run dbt models to create analytics-ready datasets
- **Monitor**: Track pipeline health and data quality

### 4. 🚀 Event-Driven Activation

Pipelines can be triggered by:

- Scheduled CronJobs for batch processing
- GCS file drops for real-time ingestion
- Pub/Sub messages for streaming data
- BigQuery table updates for data freshness

## 📊 Use Cases

### 📈 Batch Data Processing

```yaml
# Daily sales data pipeline
spec:
  ingest:
    type: cronjob
    schedule: "0 2 * * *" # Daily at 2 AM
  transform:
    models:
      - stg_daily_sales
      - fct_sales
```

### ⚡ Real-Time Event Processing

```yaml
# User activity streaming pipeline
spec:
  ingest:
    type: trigger
    name: user-activity-trigger
  transform:
    models:
      - stg_user_events
      - fct_user_sessions
```

### 🔗 Multi-Source Data Integration

```yaml
# Combine data from multiple sources
spec:
  ingest:
    type: cronjob
    schedule: "0 */6 * * *" # Every 6 hours
  transform:
    models:
      - stg_customer_data
      - stg_order_data
      - fct_customer_orders
```

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (v1.11.3+)
- kubectl configured for your cluster
- Docker or container runtime
- Access to container registry

### Installation

1. **Clone and Deploy**

   ```bash
   # Clone the repository
   git clone https://github.com/your-org/pipeline-forge.git
   cd pipeline-forge/operator

   # Deploy the operator
   make deploy IMG=your-registry/pipeline-forge-operator:latest
   ```

2. **Deploy Sample Pipelines**

   ```bash
   # Apply example configurations
   kubectl apply -k operator/config/samples/
   ```

3. **Monitor Pipeline Status**
   ```bash
   # Check pipeline health
   kubectl get staging
   kubectl describe staging user-events-staging
   ```

## 📁 Project Structure

```
pipeline-forge/
├── operator/           # Kubernetes operator (Go)
│   ├── api/           # CRD definitions
│   ├── controllers/   # Reconciliation logic
│   └── config/        # Deployment manifests
├── workloads/         # Data processing components
│   ├── ingest/        # Type-safe ingestion (Python)
│   ├── transform/     # dbt transformations
│   └── trigger/       # Event processing
└── docs/             # Documentation (planned)
```

## 🔧 Components

### Operator (`operator/`)

Kubernetes operator that manages pipeline orchestration

- [Operator Documentation](operator/README.md)

### Workloads (`workloads/`)

Specialized data processing components

- **Ingest**: Type-safe data ingestion with Pydantic
- **Transform**: dbt-based data transformation
- **Trigger**: Event processing and monitoring

## 🤝 Contributing

We welcome contributions to Pipeline Forge! Whether you're interested in:

- **Operator Development**: Kubernetes controller logic and CRDs
- **Workload Development**: Data processing components
- **Documentation**: User guides and examples
- **Testing**: End-to-end pipeline validation

Please see our contributing guidelines and development setup instructions in the respective component directories.

### 🛠️ Development Setup

```bash
# Set up development environment
cd operator
make manifests generate
make test

# Run locally
make run
```

## 📄 License

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
