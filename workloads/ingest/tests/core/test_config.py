"""Configuration validation tests for the ingest workload.

These tests demonstrate:
- Pydantic model validation
- Type safety and data validation
- Error handling for invalid configurations
- Production-ready configuration patterns
"""

import pytest
from pydantic import ValidationError

from ingest.core.config import (
    SourceConfig,
    BigQueryTargetConfig,
    Config,
    DatabaseType,
    TargetType,
    SecretProvider,
    SecretConfig,
    RuntimeParams,
)


class TestSourceConfig:
    """Test SourceConfig validation and behavior."""

    def test_valid_mysql_config(self):
        """Test valid MySQL configuration creation."""
        config = SourceConfig(
            name="test_mysql",
            type=DatabaseType.MYSQL,
            host="localhost",
            port=3306,
            username="test_user",
            password="test_password",
            database="test_db",
            ssl_required=True,
        )

        assert config.name == "test_mysql"
        assert config.type == DatabaseType.MYSQL
        assert config.host == "localhost"
        assert config.port == 3306
        assert config.ssl_required is True

    def test_valid_postgres_config(self):
        """Test valid PostgreSQL configuration creation."""
        config = SourceConfig(
            name="test_postgres",
            type=DatabaseType.POSTGRES,
            host="localhost",
            port=5432,
            username="test_user",
            password="test_password",
            database="test_db",
            schema="public",  # Use alias name
            ssl_required=False,
        )

        assert config.type == DatabaseType.POSTGRES
        assert config.db_schema == "public"  # Access via field name
        assert config.ssl_required is False

    def test_invalid_config_missing_required_fields(self):
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            SourceConfig(
                name="test",
                # Missing required fields: type, host, port, username, password, database
            )

        # Verify specific validation errors
        errors = exc_info.value.errors()
        error_fields = [error["loc"][0] for error in errors]
        assert "type" in error_fields
        assert "host" in error_fields
        assert "port" in error_fields

    def test_invalid_port_type(self):
        """Test that invalid port types raise ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            SourceConfig(
                name="test",
                type=DatabaseType.MYSQL,
                host="localhost",
                port="invalid_port",  # Should be int
                username="user",
                password="pass",
                database="db",
            )

        errors = exc_info.value.errors()
        port_error = next(error for error in errors if error["loc"][0] == "port")
        assert "int" in str(port_error["type"])

    def test_connection_string_building_mysql(self):
        """Test MySQL connection string building."""
        config = SourceConfig(
            name="test",
            type=DatabaseType.MYSQL,
            host="localhost",
            port=3306,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        conn_str = config.build_connection_string("mysql+pymysql")
        expected = "mysql+pymysql://test_user:test_password@localhost:3306/test_db"
        assert conn_str == expected

    def test_connection_string_building_postgres(self):
        """Test PostgreSQL connection string building."""
        config = SourceConfig(
            name="test",
            type=DatabaseType.POSTGRES,
            host="localhost",
            port=5432,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        conn_str = config.build_connection_string("postgresql+psycopg2")
        expected = "postgresql+psycopg2://test_user:test_password@localhost:5432/test_db"
        assert conn_str == expected

    def test_connection_string_with_default_port(self):
        """Test connection string building with default port."""
        config = SourceConfig(
            name="test",
            type=DatabaseType.MYSQL,
            host="localhost",
            port=3306,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        # Use default port (0) to test default port logic
        conn_str = config.build_connection_string("mysql+pymysql", default_port=0)
        expected = "mysql+pymysql://test_user:test_password@localhost:3306/test_db"
        assert conn_str == expected


class TestBigQueryTargetConfig:
    """Test BigQueryTargetConfig validation and behavior."""

    def test_valid_bigquery_config(self):
        """Test valid BigQuery configuration creation."""
        config = BigQueryTargetConfig(
            name="test_bq",
            type=TargetType.BIGQUERY,
            project_id="test-project",
            project_number=123456789,
            dataset_id="test_dataset",
            location="US",
        )

        assert config.name == "test_bq"
        assert config.type == TargetType.BIGQUERY
        assert config.project_id == "test-project"
        assert config.project_number == 123456789
        assert config.dataset_id == "test_dataset"
        assert config.location == "US"

    def test_invalid_project_id_format(self):
        """Test that invalid project_id format raises ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            BigQueryTargetConfig(
                name="test",
                project_id="Invalid_Project_ID",  # Invalid format
                project_number=123456789,
                dataset_id="test_dataset",
            )

        errors = exc_info.value.errors()
        project_id_error = next(error for error in errors if error["loc"][0] == "project_id")
        assert "pattern" in str(project_id_error["type"])

    def test_invalid_dataset_id_format(self):
        """Test that invalid dataset_id format raises ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            BigQueryTargetConfig(
                name="test",
                project_id="test-project",
                project_number=123456789,
                dataset_id="123invalid",  # Invalid format - starts with number
            )

        errors = exc_info.value.errors()
        dataset_error = next(error for error in errors if error["loc"][0] == "dataset_id")
        assert "pattern" in str(dataset_error["type"])


class TestSecretConfig:
    """Test SecretConfig validation and behavior."""

    def test_valid_secret_config(self):
        """Test valid secret configuration creation."""
        config = SecretConfig(
            provider=SecretProvider.GOOGLE_SECRET_MANAGER,
            name="test-secret",
            version="latest",
            secret_path="projects/123/secrets/test-secret/versions/latest",
        )

        assert config.provider == SecretProvider.GOOGLE_SECRET_MANAGER
        assert config.name == "test-secret"
        assert config.version == "latest"
        assert config.secret_path == "projects/123/secrets/test-secret/versions/latest"

    def test_secret_config_defaults(self):
        """Test secret configuration with default values."""
        config = SecretConfig(
            name="test-secret",
            # provider defaults to GOOGLE_SECRET_MANAGER
            # version defaults to "latest"
            # secret_path defaults to None
        )

        assert config.provider == SecretProvider.GOOGLE_SECRET_MANAGER
        assert config.version == "latest"
        assert config.secret_path is None


class TestRuntimeParams:
    """Test RuntimeParams validation and behavior."""

    def test_valid_runtime_params(self):
        """Test valid runtime parameters creation."""
        params = RuntimeParams(
            retry_attempts=5,
            retry_delay_seconds=60,
            chunk_size=5000,
        )

        assert params.retry_attempts == 5
        assert params.retry_delay_seconds == 60
        assert params.chunk_size == 5000

    def test_runtime_params_defaults(self):
        """Test runtime parameters with default values."""
        params = RuntimeParams()

        assert params.retry_attempts == 3  # Default value
        assert params.retry_delay_seconds == 30  # Default value
        assert params.chunk_size == 10000  # Default value

    def test_invalid_retry_attempts_range(self):
        """Test that retry_attempts outside valid range raises ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            RuntimeParams(retry_attempts=15)  # Outside 1-10 range

        errors = exc_info.value.errors()
        retry_error = next(error for error in errors if error["loc"][0] == "retry_attempts")
        assert "less_than_equal" in str(retry_error["type"])  # le=10 means less_than_equal

    def test_chunk_size_accepts_negative_values(self):
        """Test that chunk_size accepts negative values (no constraint defined)."""
        # Note: chunk_size has no minimum constraint, so negative values are allowed
        params = RuntimeParams(chunk_size=-1)
        assert params.chunk_size == -1


class TestConfigIntegration:
    """Test Config class integration and complex scenarios."""

    def test_config_get_source_config(self):
        """Test getting source configuration by environment and name."""
        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[],
            sources={
                "dev": {
                    "test_source": SourceConfig(
                        name="test_source",
                        type=DatabaseType.MYSQL,
                        host="localhost",
                        port=3306,
                        username="user",
                        password="pass",
                        database="db",
                    )
                }
            },
            targets={
                "dev": BigQueryTargetConfig(
                    name="test_target",
                    project_id="test-project",
                    project_number=123456789,
                    dataset_id="test_dataset",
                )
            },
        )

        source_config = config.get_source_config("dev", "test_source")
        assert source_config is not None
        assert source_config.name == "test_source"
        assert source_config.type == DatabaseType.MYSQL

    def test_config_get_source_config_not_found(self):
        """Test getting non-existent source configuration returns None."""
        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[],
            sources={},
            targets={},
        )

        source_config = config.get_source_config("dev", "nonexistent")
        assert source_config is None

    def test_config_extra_fields_forbidden(self):
        """Test that extra fields in Config raise ValidationError."""
        with pytest.raises(ValidationError) as exc_info:
            Config(
                version="1.0.0",
                params=RuntimeParams(),
                secrets=[],
                sources={},
                targets={},
                extra_field="not_allowed",  # This should cause validation error
            )

        errors = exc_info.value.errors()
        assert any("extra_forbidden" in str(error["type"]) for error in errors)
