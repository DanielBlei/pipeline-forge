# Pipeline Forge - Ingest Workload

A modern, type-safe data ingestion pipeline built with Python 3.13+ and Pydantic.

## Features

- **Type-Safe Configuration**: Built with Pydantic for runtime validation and type checking
- **Extensible Architecture**: Easy to add new data sources with Union types
- **Professional Standards**: Follows modern Python best practices
- **Comprehensive Logging**: Structured logging with debug support
- **CLI Interface**: Powered by Typer for excellent user experience

## Quick Start

```bash
# Install dependencies
pip install -e .

# Run with example configuration
python -m ingest.main --config example_config.yaml
```

## Configuration

The ingest workload uses Pydantic models for configuration validation. See `example_config.yaml` for a complete example.

### Supported Sources

- **MySQL**: Full MySQL support with SSL and connection pooling
- **PostgreSQL**: Native PostgreSQL support with advanced features

## Architecture

This project demonstrates modern Python development practices:

- **Pydantic Models**: Type-safe configuration and data validation
- **Union Types**: Type hints for multiple source implementations
- **Factory Pattern**: Clean source instantiation based on configuration
- **Comprehensive Error Handling**: Graceful failure with detailed logging

## Development

```bash
# Install development dependencies
pip install -e ".[dev]"

# Run tests
pytest

# Format code
ruff format .

# Lint code
ruff check .
```
