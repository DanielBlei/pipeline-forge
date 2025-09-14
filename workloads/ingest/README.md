# Pipeline Forge - Ingest Workload

A production-ready data ingestion pipeline designed for the Pipeline Forge platform.

Deployable as a standalone Docker image with comprehensive configuration management and modern Python development practices.

## Quick Start

```bash
# Install dependencies and run with `--help` flag
make install
make run-help

# Run ingestion with config files
uv run ingest --config config.yaml --catalog catalog.yaml --env dev

# Dry run (extract data but not load it into target)
uv run ingest --config config.yaml --catalog catalog.yaml --env dev --dry-run
```

## When to Use

- **Database Extraction**: MySQL, PostgreSQL → BigQuery (Future support for other targets)
- **Large Datasets**: Streaming processing with configurable chunk sizes, optimizing memory usage
- **Production Pipelines**: Built-in retry logic and error handling
- **Multi-Environment**: Dev/prod configuration management
- **Modern Python**: Type hints, async patterns, modern tooling

## How It Works

### 1. Configuration-Driven

Example configuration file:
```yaml
# Type-safe configuration with Pydantic validation
version: 1.0.0

sources:
prod:
    pipeline_forge:
      name: pipeline_forge 
      type: mysql
      host: localhost
      port: 3306
      username: pipeline_user
      password: docker-compose-mysql-password # refers to a secret name in the secrets section
      database: pipeline_forge
      ssl_required: true
  dev:
    forge_dev:
      ...
targets:
  dev:
    bigquery_target:
      project_number: 1234567890
      dataset: my_dataset
      ...
secrets:
  - name: secret-name
    provider: google_secret_manager
    ...
params:
  retry_attempts: 1
  retry_delay_seconds: 30
  chunk_size: 10000
```

Example catalog file:

```yaml
name: pipeline_forge
source: pipeline_forge
tables:
  - name: account
    replication: TRUNCATE
    columns:
      - name: id
        type: int
      - name: email
        type: string
  - name: events
    ...
```


### 2. Protocol-Based Architecture

```python
# Clean abstractions enable easy extension
source = create_source(config)  # Factory pattern
for chunk in source.extract(table, chunk_size=1000):
    target.load(chunk, write_disposition=TRUNCATE)
```

### 3. Streaming Processing

- **Memory Efficient**: Processes data in configurable chunks
- **Retry Logic**: Automatic retry with exponential backoff
- **Error Handling**: Graceful failure with detailed logging

## Key Software Engineering Practices

- **Type Safety**: Pydantic models with runtime validation
- **Protocol Design**: Clean interfaces for extensibility
- **Factory Pattern**: Extensible source/target creation
- **Comprehensive Testing**: Unit, integration, and error handling tests
- **Modern Python**: Type hints, async patterns, modern tooling

## Configuration Reference

| Field     | Type | Description                                |
| --------- | ---- | ------------------------------------------ |
| `version` | Str  | Version of the configuration file          |
| `sources` | Dict | Environment-specific source configs        |
| `targets` | Dict | Environment-specific target configurations |
| `secrets` | List | Secret references (access during runtime)  |
| `params`  | Dict | Runtime parameters (access during runtime) |

## Development

```bash
make install     # Install dependencies
make check       # Run linters (Ruff and MyPy)
make fix         # Fix linting errors
make test        # Run tests
make test-coverage  # Run tests with coverage report and open HTML report in browser
```
