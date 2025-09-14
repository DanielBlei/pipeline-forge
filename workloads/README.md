# Pipeline Forge - Workloads

This directory contains the data processing components for Pipeline Forge. All workloads are packaged as Docker images for easy deployment and scaling.

## Available Workloads

- **[Ingest](./ingest/README.md)** - Data ingestion from databases to BigQuery (Python)
- **[Transform](./transform/README.md)** - dbt-based data transformation
- **[Trigger](./trigger/README.md)** - Event-driven pipeline activation (Go)

## Architecture Overview

All workloads follow consistent patterns:

- **Containerized Deployment** - Packaged as Docker images for easy deployment
- **Standalone Operation** - Each workload runs independently with its own configuration
- **Type-safe Configuration** - Runtime validation and environment-specific configs
- **Comprehensive Testing** - Unit, integration, and end-to-end test coverage
- **Modern Development Practices** - Linting, formatting, and CI/CD ready

_For detailed information about each workload, see their individual README files._
