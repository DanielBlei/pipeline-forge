# Pipeline Forge 🔥

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.11+-blue.svg)](https://kubernetes.io/)
[![Python](https://img.shields.io/badge/Python-3.13+-blue.svg)](https://python.org/)
[![dbt](https://img.shields.io/badge/dbt-Core-blue.svg)](https://www.getdbt.com/)

> A Kubernetes-native platform for building modern, declarative data pipelines with clear boundaries between ingestion and transformation.

## 🚀 Quick Navigation

- **[🎛️ Kubernetes Operator](operator/README.md)** - Go-based CRD management and pipeline orchestration
- **[📥 Ingest Workload](workloads/ingest/README.md)** - Type-safe data ingestion (Python)
- **[🔄 Transform Workload](workloads/transform/README.md)** - dbt-based data transformation
- **[⚡ Trigger Workload](workloads/trigger/README.md)** - Event-driven pipeline activation (Go)
- **[📋 Examples](docs/examples.md)** - Comprehensive YAML examples and use cases

---

## 🎯 What is Pipeline Forge?

A complete solution for orchestrating data pipelines in Kubernetes environments. Combines a powerful Kubernetes operator with specialized workloads to provide a declarative, event-driven approach to data pipeline management.

### Key Benefits

- **Unified Pipeline Lifecycle** - Connect ingestion with transformation in a single application lifecycle
- **Native Kubernetes Resources** - Each step runs on 100% native K8s resources
- **Event-Driven Orchestration** - React to file drops, Pub/Sub messages, and BigQuery updates
- **Built-in Observability** - Comprehensive status tracking and monitoring

## 🏗️ Architecture Overview

Pipeline Forge consists of two main components:

### 🎛️ [Kubernetes Operator](operator/README.md)

**Go-based CRD management and pipeline orchestration**

- Custom Resource Definitions (CRDs) for pipeline definition
- Automatic reconciliation and lifecycle management
- RBAC integration and resource management
- Event-driven trigger management

### 🔧 [Specialized Workloads](workloads/README.md)

**Production-ready data processing components**

- **[Ingest](workloads/ingest/README.md)** - Type-safe data ingestion from MySQL, PostgreSQL to BigQuery
- **[Transform](workloads/transform/README.md)** - dbt-based data transformation with version control
- **[Trigger](workloads/trigger/README.md)** - Event processing for GCS, Pub/Sub, and BigQuery

## 🛠️ Technology Stack

| Component     | Technology                    | Purpose                                   |
| ------------- | ----------------------------- | ----------------------------------------- |
| **Operator**  | Go, Kubernetes, Kubebuilder   | Pipeline orchestration and CRD management |
| **Ingest**    | Python 3.13+, Pydantic, Typer | Type-safe data ingestion with validation  |
| **Transform** | dbt Core, BigQuery            | Data transformation and analytics         |
| **Triggers**  | Go, Google Cloud APIs         | Event-driven pipeline activation          |

## 🚀 Quick Start

```bash
git clone https://github.com/DanielBlei/pipeline-forge.git

# Run the operator (safety check ensures you're on kind/minikube cluster)
make run-operator

# Deploy k8s samples resources
make apply-samples
```

## 📁 Project Structure

```
pipeline-forge/
├── operator/           # Kubernetes operator (Go)
├── workloads/          # Data processing components
│   ├── ingest/        # Type-safe ingestion (Python)
│   ├── transform/     # dbt transformations
│   └── trigger/       # Event processing (Go)
└── docs/              # Documentation
```

## 🤝 Contributing

We welcome contributions! See individual component READMEs for development setup and guidelines.

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
