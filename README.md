# Pipeline Forge 🔥

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.11+-blue.svg)](https://kubernetes.io/)
[![Python](https://img.shields.io/badge/Python-3.13+-blue.svg)](https://python.org/)
[![dbt](https://img.shields.io/badge/dbt-Core-blue.svg)](https://www.getdbt.com/)

> A Kubernetes-native platform for building modern, declarative data pipelines with clear boundaries between ingestion and transformation.

## 📋 Quick Navigation

- [Overview](#-overview)
- [Key Features](#-key-features)
- [Quick Start](#-quick-start)
- [Architecture](#️-architecture)
- [Documentation](#-documentation)
- [Contributing](#-contributing)

---

## 🚀 Overview

Pipeline Forge is a complete solution for orchestrating data pipelines in Kubernetes environments. It combines a powerful Kubernetes operator with specialized workloads to provide a declarative, event-driven approach to data pipeline management.

### 🎯 The Problem

Modern data teams face challenges with:

- **Complex Orchestration**: Managing dependencies between data ingestion and transformation
- **Event-Driven Requirements**: Responding to file drops, database changes, and streaming events
- **Infrastructure Complexity**: Deploying and scaling data processing workloads
- **Observability Gaps**: Tracking pipeline health and data lineage
- **Team Coordination**: Coordinating between data engineering and platform teams
- **Resilience and Lifecycle Management**: Ensuring each pipeline step is connected in a clear lifecycle—if one step fails, others don't attempt to run, preventing cascading errors and maintaining robust execution

### 💡 The Solution

Pipeline Forge provides a Kubernetes-native platform that:

- **Declaratively Orchestrates** data pipelines using Custom Resource Definitions (CRDs)
- **Flexible Ingestion** supports both event-driven (e.g., GCS file drops, Pub/Sub messages, BigQuery updates) and scheduled (CronJob-based) pipeline execution
- **Clear Separation** between data ingestion and transformation phases
- **Built-in Observability** with comprehensive status tracking and monitoring
- **GitOps Ready** configuration that fits modern deployment practices

## ✨ Key Features

- **Unified Pipeline Lifecycle**: Connect ingestion with staging models in a single application lifecycle - if ingestion fails, the entire staging fails, preventing orphaned transformations
- **Native Kubernetes Resources**: Each step runs on 100% native K8s resources (Transform → Job, Ingest → CronJob/Job/Trigger)
- **Event-Driven Orchestration**: React to file drops, Pub/Sub messages, and BigQuery updates with intelligent retry policies
- **Built-in Observability**: Comprehensive status tracking with detailed execution history and failure analysis
- **Flexible Ingestion**: Reference existing CronJobs during Ingestion or create new ones as needed with full type safety and managed by the operator
- **Custom Image Support**: Use your own image for each step, or use pre-built Docker images from the Pipeline Forge repository

### ⚡ Event-Driven Orchestration

- **GCS Triggers**: Monitor bucket changes and trigger pipelines
- **Pub/Sub Triggers**: React to real-time messages with optional filtering
- **BigQuery Triggers**: Watch for table updates and data freshness
- **Retry & Cooldown**: Configurable retry policies with intelligent intervals

### ☸️ Kubernetes-Native

- **CRD-Based**: Native Kubernetes resources for pipeline definition
- **RBAC Integration**: Fine-grained access control for teams
- **Resource Management**: CPU, memory, and storage allocation
- **Independent Scaling**: Each step scales independently as native K8s resources

### 📊 Built-in Observability

- **Rich Status Tracking**: Comprehensive pipeline health monitoring with detailed execution history
- **Lifecycle Management**: Real-time phase tracking (Pending, Running, Completed, Failed)
- **Execution Insights**: Track attempt counts, success/failure rates, and timing metrics
- **Failure Analysis**: Detailed error messages and retry attempt tracking

### 🔄 Example

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: user-events-pipeline
  namespace: staging-events
spec:
  ingest:
    mode: reference
    type: trigger
    name: user-events-trigger
  transform:
    name: user-events-transform
    project: analytics
    target: prod
    image: gcr.io/org/dbt-core:latest
    models:
      - stg_user_events
```

📖 **[View comprehensive examples →](docs/examples.md)**

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

## 🏗️ Architecture

Pipeline Forge consists of two main components that work together seamlessly:

### 🎛️ Kubernetes Operator

The operator manages the lifecycle of data pipeline stages through:

- **Staging Resources**: Complete pipeline steps that coordinate ingestion and transformation
- **Trigger Resources**: Event-driven activation for pipelines
- **Ingestion Management**: Supports ingestion via both event-driven triggers and CronJobs
- **Automatic Reconciliation**: Ensures pipeline state matches desired configuration

### 🔧 Specialized Workloads

Pre-built, production-ready data processing components:

- **Ingest Workloads**: Type-safe data ingestion from MySQL, PostgreSQL, and more
- **Transform Workloads**: dbt-based data transformation with version control
- **Trigger Workloads**: Event processing for GCS, Pub/Sub, and BigQuery

## 📚 Documentation

- **[📋 Examples](docs/examples.md)** - Comprehensive YAML examples and use cases
- **[🎛️ Operator Guide](operator/README.md)** - Detailed operator documentation
- **[🔧 Workloads](workloads/README.md)** - Data processing components

## 🛠️ Technology Stack

| Component          | Technology                    | Purpose                                              |
| ------------------ | ----------------------------- | ---------------------------------------------------- |
| **Operator**       | Go, Kubernetes, Kubebuilder   | Pipeline orchestration and CRD management            |
| **Ingest**         | Python 3.13+, Pydantic, Typer | Type-safe data ingestion with validation             |
| **Transform**      | dbt Core, BigQuery            | Data transformation and analytics                    |
| **Triggers**       | GCS, Pub/Sub, BigQuery APIs   | Event-driven pipeline activation with retry policies |
| **Infrastructure** | Kubernetes, Docker            | Container orchestration and deployment               |

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
└── docs/              # Documentation
    └── examples.md    # Comprehensive examples
```

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

```

```
