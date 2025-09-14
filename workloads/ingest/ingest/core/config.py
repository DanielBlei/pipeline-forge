"""Configuration classes for ingestion"""

from typing_extensions import Union
from pydantic import BaseModel, Field, ConfigDict
from typing import Dict, List, Optional
from enum import Enum

from ingest.helpers.secret_handler import get_gcloud_secret


class DatabaseType(str, Enum):
    MYSQL = "mysql"
    POSTGRES = "postgres"


class TargetType(str, Enum):
    BIGQUERY = "bigquery"


class SourceConfig(BaseModel):
    name: str
    type: DatabaseType
    host: str
    port: int
    username: str
    password: str = Field(..., description="Either refer to a secret name in the Secrets section")
    database: str
    db_schema: Optional[str] = Field(None, alias="schema")
    ssl_required: Optional[bool] = False

    def build_connection_string(self, dialect: str, default_port: int = 0) -> str:
        """Build database connection string from config.

        Args:
            dialect: SQLAlchemy dialect (e.g., 'mysql+pymysql', 'postgresql+psycopg2')
            default_port: Default port for the database (will use config.port if None)

        Returns:
            SQLAlchemy connection string
        """
        port = default_port if default_port != 0 else self.port
        return f"{dialect}://{self.username}:{self.password}@{self.host}:{port}/{self.database}"


class BigQueryTargetConfig(BaseModel):
    name: str
    type: TargetType = TargetType.BIGQUERY
    project_id: str = Field(
        ...,
        pattern=r"^[a-z][a-z0-9\-]{4,28}[a-z0-9]$",
        description="BigQuery project_id must be 6-30 characters, lowercase letters, digits or hyphens, start with a letter, end with letter or digit.",
    )
    project_number: int = Field(
        ...,
        description="BigQuery project_number must be a valid project number.",
    )
    dataset_id: str = Field(
        ...,
        pattern=r"^[A-Za-z_][A-Za-z0-9_]{0,1023}$",
        description="BigQuery dataset_id must be 1-1024 characters, start with a letter or underscore, contain only letters, numbers, or underscores.",
    )
    location: Optional[str] = None
    service_account: Optional[str] = None


TargetTypes = Union["BigQueryTargetConfig"]


class RuntimeParams(BaseModel):
    retry_attempts: int = Field(ge=1, le=10, default=3)
    retry_delay_seconds: int = Field(ge=1, le=3600, default=30)
    chunk_size: int = Field(default=10000)


class SecretProvider(str, Enum):
    GOOGLE_SECRET_MANAGER = "gcloud"
    # AWS_SECRET_MANAGER = "aws"
    # AZURE_KEY_VAULT = "azure"


class SecretConfig(BaseModel):
    provider: SecretProvider = Field(
        default=SecretProvider.GOOGLE_SECRET_MANAGER, description="Secret Manager Provider"
    )
    name: str = Field(description="Secret name")
    version: Optional[str] = Field(default="latest", description="Secret version")
    secret_path: Optional[str] = Field(
        default=None,
        description="Secret path in the Gcloud Secret Manager(e.g projects/pipeline-forge/secrets/postgres-password)",
    )


class Config(BaseModel):
    version: str
    params: RuntimeParams
    secrets: List[SecretConfig]
    sources: Dict[str, Dict[str, SourceConfig]]  # environment -> source_name -> SourceConfig
    targets: Dict[str, TargetTypes]  # enviroument -> TargetType

    model_config = ConfigDict(extra="forbid", validate_assignment=True, str_strip_whitespace=True)

    def get_source_config(self, environment: str, source_name: str) -> Optional[SourceConfig]:
        """Get a specific source configuration by environment and name"""
        if environment not in self.sources:
            return None

        return self.sources[environment].get(source_name)

    def get_target_config(self, environment: str) -> Optional[TargetTypes]:
        """Get a specific target configuration by environment and name"""
        if environment not in self.targets:
            return None

        return self.targets.get(environment)

    def get_gcloud_secret_value(self, secret_name: str, environment: str, version: Optional[str] = "latest") -> str:
        """Get a secret from the Gcloud Secret Manager"""
        if secret_name not in [secret.name for secret in self.secrets]:
            raise ValueError(f"Secret {secret_name} not found in secrets")

        secret_config = next(secret for secret in self.secrets if secret.name == secret_name)
        if secret_config.provider != SecretProvider.GOOGLE_SECRET_MANAGER:
            raise ValueError(f"Secret {secret_name} is not a Gcloud Secret")

        if secret_config.secret_path is None:
            target_config = self.targets.get(environment)
            if target_config is None:
                raise ValueError(f"Target config not found for environment {environment}")

            secret_path = f"projects/{target_config.project_number}/secrets/{secret_name}/versions/{version}"
        else:
            secret_path = secret_config.secret_path

        return get_gcloud_secret(secret_path)
