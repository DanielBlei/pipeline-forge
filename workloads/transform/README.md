# Pipeline Forge - Transform Workload

A starting point for data transformation using DBT, designed for the Pipeline Forge sandbox environment.

Basic DBT project structure with service account integration for BigQuery access.

**Infrastructure Reference**: See the [Infrastructure README](../../infrastructure/README.md) for sandbox environment setup and service account configuration.

## Quick Start

```bash
# Run transformations
make run

# Run tests
make test

# Build Docker image
make build-docker

# Run Transform Docker image
docker run --rm transform dbt help

# Show all available commands
make help
```

## How It Works

### 1. DBT Project Structure

The project follows DBT best practices with organized model layers:

```
pipeline_forge/
├── models/
│   └── pipeline_forge/
│       ├── stg_pipeline_forge__accounts.sql
│       ├── stg_pipeline_forge__events.sql
│       ├── stg_pipeline_forge__user_events.sql
│       └── pipeline_forge_schema.yml
├── dbt_project.yml
├── profiles.yml
└── packages.yml
```

### 2. Service Account Integration

Secure BigQuery access using dedicated service account:

```yaml
# profiles.yml
pipeline_forge:
  outputs:
    dev:
      type: bigquery
      project: pipeline-forge
      dataset: pipeline_forge_transform
      impersonate_service_account: transform-workload@pipeline-forge.iam.gserviceaccount.com
      method: oauth
      location: us-central1
      threads: 1
      job_execution_timeout_seconds: 300
      job_retries: 1
      priority: interactive
  target: dev
```

### 3. Staging Models

Clean, consistent staging layer with standardized naming:

- **Naming Convention**: `stg_{source}__{table}`
- **Materialization**: Views for fast iteration
- **Schema**: Organized in `stg_pipeline_forge` schema
- **Documentation**: Comprehensive schema documentation in YAML

## Key Features

- **Service Account Security**: Dedicated IAM service account for BigQuery access
- **Basic DBT Structure**: Standard project layout with staging models
- **Sandbox Ready**: Configured for Pipeline Forge sandbox environment

## Configuration Reference

| File                        | Purpose               | Key Settings                               |
| --------------------------- | --------------------- | ------------------------------------------ |
| `dbt_project.yml`           | Project configuration | Model paths, materialization settings      |
| `profiles.yml`              | Connection settings   | BigQuery project, service account, dataset |
| `packages.yml`              | DBT packages          | dbt_utils, codegen dependencies            |
| `pipeline_forge_schema.yml` | Model documentation   | Column descriptions, tests, relationships  |
