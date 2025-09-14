"""Factory function tests for the ingest workload.

These tests demonstrate:
- Factory pattern implementation
- Dependency injection and validation
- Error handling in factory functions
- Mocking external dependencies
"""

import pytest
from unittest.mock import Mock, patch

from ingest.sources import create_source
from ingest.targets import create_target
from ingest.core.config import SourceConfig, BigQueryTargetConfig, DatabaseType, TargetType


class TestSourceFactory:
    """Test create_source factory function."""

    @patch("ingest.sources.MySQLSource")
    def test_create_mysql_source_success(self, mock_mysql_source):
        """Test successful MySQL source creation."""
        # Setup
        mock_source_instance = Mock()
        mock_source_instance.validate_connection.return_value = True
        mock_mysql_source.return_value = mock_source_instance

        config = SourceConfig(
            name="test_mysql",
            type=DatabaseType.MYSQL,
            host="localhost",
            port=3306,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        # Execute
        result = create_source(config)

        # Verify
        mock_mysql_source.assert_called_once_with(config)
        mock_source_instance.validate_connection.assert_called_once()
        assert result == mock_source_instance

    @patch("ingest.sources.PostgresSource")
    def test_create_postgres_source_success(self, mock_postgres_source):
        """Test successful PostgreSQL source creation."""
        # Setup
        mock_source_instance = Mock()
        mock_source_instance.validate_connection.return_value = True
        mock_postgres_source.return_value = mock_source_instance

        config = SourceConfig(
            name="test_postgres",
            type=DatabaseType.POSTGRES,
            host="localhost",
            port=5432,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        # Execute
        result = create_source(config)

        # Verify
        mock_postgres_source.assert_called_once_with(config)
        mock_source_instance.validate_connection.assert_called_once()
        assert result == mock_source_instance

    def test_create_source_unsupported_type(self):
        """Test that unsupported source types raise ValueError."""
        # Create a mock config that bypasses Pydantic validation
        # but has the right structure for the factory function
        config = Mock()
        config.type.value = "unsupported_type"

        with pytest.raises(ValueError, match="Unsupported source type"):
            create_source(config)

    @patch("ingest.sources.MySQLSource")
    def test_create_source_connection_validation_failure(self, mock_mysql_source):
        """Test that connection validation failure raises ValueError."""
        # Setup
        mock_source_instance = Mock()
        mock_source_instance.validate_connection.return_value = False
        mock_mysql_source.return_value = mock_source_instance

        config = SourceConfig(
            name="test_mysql",
            type=DatabaseType.MYSQL,
            host="localhost",
            port=3306,
            username="test_user",
            password="test_password",
            database="test_db",
        )

        # Execute & Verify
        with pytest.raises(ValueError, match="Failed to validate source connection"):
            create_source(config)


class TestTargetFactory:
    """Test create_target factory function."""

    @patch("ingest.targets.BigQueryTarget")
    def test_create_bigquery_target_success(self, mock_bigquery_target):
        """Test successful BigQuery target creation."""
        # Setup
        mock_target_instance = Mock()
        mock_target_instance.validate_connection.return_value = True
        mock_bigquery_target.return_value = mock_target_instance

        config = BigQueryTargetConfig(
            name="test_bq",
            type=TargetType.BIGQUERY,
            project_id="test-project",
            project_number=123456789,
            dataset_id="test_dataset",
        )

        # Execute
        result = create_target(config)

        # Verify
        mock_bigquery_target.assert_called_once_with(config)
        mock_target_instance.validate_connection.assert_called_once()
        assert result == mock_target_instance

    def test_create_target_unsupported_type(self):
        """Test that unsupported target types raise ValueError."""
        # Create a real BigQueryTargetConfig but manually set invalid type
        config = BigQueryTargetConfig(
            name="test_unsupported",
            type=TargetType.BIGQUERY,  # Start with valid type
            project_id="test-project",
            project_number=123456789,
            dataset_id="test_dataset",
        )
        # Manually change the type to simulate unsupported type
        config.type = "unsupported_type"

        with pytest.raises(ValueError, match="Unsupported target type"):
            create_target(config)

    @patch("ingest.targets.BigQueryTarget")
    def test_create_target_connection_validation_failure(self, mock_bigquery_target):
        """Test that connection validation failure raises ValueError."""
        # Setup
        mock_target_instance = Mock()
        mock_target_instance.validate_connection.return_value = False
        mock_bigquery_target.return_value = mock_target_instance

        config = BigQueryTargetConfig(
            name="test_bq",
            type=TargetType.BIGQUERY,
            project_id="test-project",
            project_number=123456789,
            dataset_id="test_dataset",
        )

        # Execute & Verify
        with pytest.raises(ValueError, match="Failed to validate target connection"):
            create_target(config)
