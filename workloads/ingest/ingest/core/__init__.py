"""Core package for pipeline-forge ingest functionality."""

from .config import Config, SourceConfig, RuntimeParams, SecretsConfig, DatabaseType, BigQueryTarget

from .catalog import Catalog, Table, Column, ReplicationType

__all__ = [
    # Config models
    "Config",
    "SourceConfig",
    "RuntimeParams",
    "SecretsConfig",
    "DatabaseType",
    "BigQueryTarget",
    # Catalog models
    "Catalog",
    "Table",
    "Column",
    "ReplicationType",
]
