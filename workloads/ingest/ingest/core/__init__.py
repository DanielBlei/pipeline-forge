"""Core package for pipeline-forge ingest functionality."""

from .config import Config, SourceConfig, RuntimeParams, DatabaseType, BigQueryTargetConfig

from .catalog import Catalog, Table, Column, ReplicationType

__all__ = [
    # Config models
    "Config",
    "SourceConfig",
    "RuntimeParams",
    "DatabaseType",
    "BigQueryTargetConfig",
    # Catalog models
    "Catalog",
    "Table",
    "Column",
    "ReplicationType",
]
