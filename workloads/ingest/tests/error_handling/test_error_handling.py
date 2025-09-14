"""Error handling tests for the ingest workload.

These tests demonstrate:
- Graceful failure handling
- Production-ready error scenarios
- Resilience patterns
- Error recovery strategies
"""

import pytest
from unittest.mock import Mock, patch

from ingest.core.config import Config, RuntimeParams, SecretConfig, SecretProvider
from ingest.core.catalog import Catalog, Table, Column, ReplicationType
from ingest.helpers.secret_handler import get_gcloud_secret


class TestSecretErrorHandling:
    """Test secret handling error scenarios."""

    @patch("ingest.helpers.secret_handler.secretmanager.SecretManagerServiceClient")
    def test_secret_not_found_error(self, mock_client_class):
        """Test handling of secret not found error."""
        # Setup
        mock_client = Mock()
        mock_client.access_secret_version.side_effect = Exception("Secret not found")
        mock_client_class.return_value = mock_client

        # Execute & Verify
        with pytest.raises(ValueError, match="Failed to get secret"):
            get_gcloud_secret("projects/123/secrets/nonexistent/versions/latest")

    @patch("ingest.helpers.secret_handler.secretmanager.SecretManagerServiceClient")
    def test_secret_empty_value_error(self, mock_client_class):
        """Test handling of empty secret value."""
        # Setup
        mock_client = Mock()
        mock_secret = Mock()
        mock_secret.payload.data.decode.return_value = None
        mock_client.access_secret_version.return_value = mock_secret
        mock_client_class.return_value = mock_client

        # Execute & Verify
        with pytest.raises(ValueError, match="Secret .* is empty"):
            get_gcloud_secret("projects/123/secrets/empty-secret/versions/latest")

    def test_config_secret_not_found_in_secrets_list(self):
        """Test that config raises error when secret not found in secrets list."""
        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[SecretConfig(name="existing-secret", provider=SecretProvider.GOOGLE_SECRET_MANAGER)],
            sources={},
            targets={},
        )

        # Execute & Verify
        with pytest.raises(ValueError, match="Secret nonexistent-secret not found in secrets"):
            config.get_gcloud_secret_value("nonexistent-secret", "dev")


class TestCatalogErrorHandling:
    """Test catalog error handling scenarios."""

    def test_catalog_get_table_not_found(self):
        """Test that getting non-existent table raises ValueError."""
        catalog = Catalog(
            name="test_catalog",
            source="test_source",
            tables=[
                Table(
                    name="existing_table",
                    replication=ReplicationType.TRUNCATE,
                    columns=[Column(name="id", type="int")],
                )
            ],
        )

        # Execute & Verify
        with pytest.raises(ValueError, match="Table nonexistent_table not found in catalog"):
            catalog.get_table("nonexistent_table")

    def test_catalog_empty_tables_list(self):
        """Test catalog behavior with empty tables list."""
        catalog = Catalog(
            name="test_catalog",
            source="test_source",
            tables=[],
        )

        # Execute
        tables = catalog.get_tables_by_source("test_source")
        sources = catalog.get_sources()

        # Verify
        assert len(tables) == 0
        assert len(sources) == 0


class TestConfigErrorHandling:
    """Test configuration error handling scenarios."""

    def test_config_get_source_config_not_found(self):
        """Test that getting non-existent source config returns None."""
        # Create a proper Config instance
        from ingest.core.config import SourceConfig, DatabaseType

        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[],
            sources={
                "dev": {
                    "existing_source": SourceConfig(
                        name="existing_source",
                        type=DatabaseType.MYSQL,
                        host="localhost",
                        port=3306,
                        username="user",
                        password="pass",
                        database="db",
                    )
                }
            },
            targets={},
        )

        # Execute
        result = config.get_source_config("dev", "nonexistent_source")

        # Verify
        assert result is None

    def test_config_get_source_config_wrong_environment(self):
        """Test that getting source config from wrong environment returns None."""
        from ingest.core.config import SourceConfig, DatabaseType

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
            targets={},
        )

        # Execute
        result = config.get_source_config("prod", "test_source")

        # Verify
        assert result is None

    def test_config_get_target_config_not_found(self):
        """Test that getting non-existent target config returns None."""
        from ingest.core.config import BigQueryTargetConfig, TargetType

        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[],
            sources={},
            targets={
                "dev": BigQueryTargetConfig(
                    name="test_target",
                    type=TargetType.BIGQUERY,
                    project_id="test-project",
                    project_number=123456789,
                    dataset_id="test_dataset",
                )
            },
        )

        # Execute
        result = config.get_target_config("prod")

        # Verify
        assert result is None


class TestRuntimeParamsErrorHandling:
    """Test runtime parameters error handling."""

    def test_runtime_params_retry_attempts_too_high(self):
        """Test that retry_attempts above maximum raises ValidationError."""
        with pytest.raises(Exception):  # ValidationError from Pydantic
            RuntimeParams(retry_attempts=15)  # Above max of 10

    def test_runtime_params_retry_attempts_too_low(self):
        """Test that retry_attempts below minimum raises ValidationError."""
        with pytest.raises(Exception):  # ValidationError from Pydantic
            RuntimeParams(retry_attempts=0)  # Below min of 1

    def test_runtime_params_retry_delay_too_high(self):
        """Test that retry_delay_seconds above maximum raises ValidationError."""
        with pytest.raises(Exception):  # ValidationError from Pydantic
            RuntimeParams(retry_delay_seconds=4000)  # Above max of 3600

    def test_runtime_params_retry_delay_too_low(self):
        """Test that retry_delay_seconds below minimum raises ValidationError."""
        with pytest.raises(Exception):  # ValidationError from Pydantic
            RuntimeParams(retry_delay_seconds=0)  # Below min of 1


class TestTableErrorHandling:
    """Test table error handling scenarios."""

    def test_table_invalid_replication_type(self):
        """Test that invalid replication type raises ValidationError."""
        with pytest.raises(Exception):  # ValidationError from Pydantic
            Table(
                name="test_table",
                replication="INVALID_TYPE",  # Invalid replication type
                columns=[Column(name="id", type="int")],
            )

    def test_table_empty_columns_list(self):
        """Test table with empty columns list (should be allowed)."""
        table = Table(
            name="test_table",
            replication=ReplicationType.TRUNCATE,
            columns=[],  # Empty columns list
        )

        # Verify - should not raise exception
        assert len(table.columns) == 0
        assert table.name == "test_table"
