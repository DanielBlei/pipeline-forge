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
- **[🛠️ Development Environment](dev/README.md)** - Local development setup and database provisioning
- **[☁️ Infrastructure](infrastructure/README.md)** - Cloud infrastructure automation and deployment
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
- **[Transform](workloads/transform/README.md)** - dbt-based data transformation
- **[Trigger](workloads/trigger/README.md)** - Event processing for GCS, Pub/Sub, and BigQuery

## 🛠️ Technology Stack

| Component           | Technology                    | Purpose                                   |
| ------------------- | ----------------------------- | ----------------------------------------- |
| **Operator**        | Go, Kubernetes, Kubebuilder   | Pipeline orchestration and CRD management |
| **Ingest**          | Python 3.13+, Pydantic, Typer | Type-safe data ingestion with validation  |
| **Transform**       | dbt Core, BigQuery            | Data transformation and analytics         |
| **Triggers**        | Go, Google Cloud APIs         | Event-driven pipeline activation          |
| **Dev Environment** | Docker Compose, SQL           | Local development and testing             |
| **Infrastructure**  | Terraform, GCP                | Cloud infrastructure automation           |

## 📊 Project Overview

### 📁 Structure

```
pipeline-forge/
├── operator/           # Kubernetes operator (Go)
├── workloads/          # Data processing components
│   ├── ingest/        # Type-safe ingestion (Python)
│   ├── transform/     # dbt transformations
│   └── trigger/       # Event processing (Go)
├── dev/               # Development environment setup
├── infrastructure/    # Cloud infrastructure automation
└── docs/              # Documentation
```

### 🚧 Status

**Current State**: Work in Progress

| Component                      | Status                | Description                                     |
| ------------------------------ | --------------------- | ----------------------------------------------- |
| **🎛️ Operator API**            | ⚡ **Functional**     | CRD definitions and API contracts               |
| **🎛️ Operator Reconciliation** | 🚧 **In Development** | Pipeline orchestration and lifecycle management |
| **📥 Ingest Workload**         | ⚡ **Functional**     | Type-safe data ingestion (Python)    |
| **🔄 Transform Workload**      | ⚡ **Functional**     | dbt-core data transformation                   |
| **⚡ Trigger Workload**        | 🚧 **In Development** | Event-driven pipeline activation (Go)           |

## 🧪 Experimental Branches

This project encourages experimental branches (prefixed with `lab/*`), current experiments include:

- **[`lab/jenkins-k8s`](../../tree/lab/jenkins-k8s/jenkins)** - Jenkins CI/CD evaluation with Kubernetes-native deployment, Configuration as Code (JCasC), and Job DSL seed patterns

These experiments help evaluate different approaches and technologies that may or may not be integrated into the main project.

## 📄 License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please open an issue first to discuss any changes before submitting a pull request.

## ⚠️ Disclaimer

This is a personal open-source project, developed independently on my own time and equipment.  
It is **not affiliated with, endorsed by, or representing my employer**.
