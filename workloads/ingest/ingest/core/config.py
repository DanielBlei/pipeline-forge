""" Configuration classes for ingestion """
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
    targets: Dict[str, Dict[str, BigQueryTarget]]

    model_config = ConfigDict(
            extra='forbid',
            validate_assignment=True,
            str_strip_whitespace=True
        )

    def get_source_config(self, environment: str, source_name: str) -> Optional[SourceConfig]:
        """Get a specific source configuration by environment and name"""
        if environment not in self.sources:
            return None
        
        return self.sources.get(environment).get(source_name)


    def get_target_config(self, environment: str, target_name: str) -> Optional[BigQueryTarget]:
        """Get a specific target configuration by environment and name"""
        if environment not in self.targets:
            return None
        
        return self.targets.get(environment).get(target_name)

