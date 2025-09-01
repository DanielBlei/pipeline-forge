"""Core package for pipeline-forge ingest functionality."""

from .config import (   
    Config,
    SourceConfig,
    TargetConfig,
    ConnectionConfig,
    RuntimeParams,
    SecretsConfig,
    DatabaseType
)

from .catalog import (
    Catalog,
    Table,
    Column,
    ReplicationType
)

__all__ = [
    # Config models
    "Config",
    "SourceConfig", 
    "TargetConfig",
    "ConnectionConfig",
    "RuntimeParams",
    "SecretsConfig",
    "DatabaseType",

    # Catalog models
    "Catalog",
    "Table",
    "Column",
    "ReplicationType",

]