"""Configuration classes for ingestion"""

from pydantic import BaseModel, Field, ConfigDict
from typing import Dict, List, Optional
from enum import Enum


class DatabaseType(str, Enum):
    MYSQL = "mysql"
    POSTGRES = "postgres"
    BIGQUERY = "bigquery"


class ConnectionConfig(BaseModel):
    host: str
    port: int
    username: str
    password: str  # Will be resolved from env or secret manager
    database: str
    db_schema: Optional[str] = Field(None, alias="schema")
    ssl_required: Optional[bool] = False

    def build_connection_string(self, dialect: str, default_port: int = None) -> str:
        """Build database connection string from config.

        Args:
            dialect: SQLAlchemy dialect (e.g., 'mysql+pymysql', 'postgresql+psycopg2')
            default_port: Default port for the database (will use config.port if None)

        Returns:
            SQLAlchemy connection string
        """
        port = default_port if default_port is not None else self.port
        return f"{dialect}://{self.username}:{self.password}@{self.host}:{port}/{self.database}"


class SourceConfig(BaseModel):
    name: str
    type: DatabaseType
    connection: ConnectionConfig
    ssl_required: bool = False


class BigQueryTarget(BaseModel):
    name: str
    type: str = "bigquery"
    project_id: str
    dataset: str
    location: Optional[str] = None
    service_account: str


class RuntimeParams(BaseModel):
    retry_attempts: int = Field(ge=1, le=10, default=3)
    retry_delay_seconds: int = Field(ge=1, le=3600, default=30)
    chunk_size: int = Field(default=10000)


class TargetConfig(BaseModel):
    name: str
    type: DatabaseType
    connection: ConnectionConfig
    ssl_required: bool = False


class SecretConfig(BaseModel):
    name: str
    path: str


class SecretsConfig(BaseModel):
    provider: str
    secrets: List[SecretConfig]


class Config(BaseModel):
    version: str
    params: RuntimeParams
    secrets: SecretsConfig
    sources: Dict[str, Dict[str, SourceConfig]]  # environment -> source_name -> SourceConfig
    targets: Dict[str, BigQueryTarget]

    model_config = ConfigDict(extra="forbid", validate_assignment=True, str_strip_whitespace=True)

    def get_source_config(self, environment: str, source_name: str) -> Optional[SourceConfig]:
        """Get a specific source configuration by environment and name"""
        if environment not in self.sources:
            return None

        return self.sources.get(environment).get(source_name)

    def get_target_config(self, environment: str) -> Optional[BigQueryTarget]:
        """Get a specific target configuration by environment and name"""
        if environment not in self.targets:
            return None

        return self.targets.get(environment)
