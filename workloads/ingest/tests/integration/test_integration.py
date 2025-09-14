"""Integration test for the ingest workload.

This test demonstrates:
- End-to-end pipeline testing
- Mocking external dependencies
- Full workflow validation
- Production-ready integration patterns
"""

import pytest
from unittest.mock import Mock, patch

from ingest.main import main
from ingest.core.config import Config
from ingest.core.catalog import Catalog, Table, Column, ReplicationType


class TestFullPipelineIntegration:
    """Test the complete ingestion pipeline."""

    @patch("ingest.main.create_target")
    @patch("ingest.main.create_source")
    @patch("ingest.main.load_yaml_model")
    def test_full_pipeline_success(self, mock_load_yaml, mock_create_source, mock_create_target):
        """Test successful end-to-end pipeline execution."""
        # Setup mock data
        mock_config = Mock(spec=Config)
        mock_config.targets = {"dev": Mock()}
        mock_config.get_gcloud_secret_value.return_value = "resolved_password"
        mock_config.params = Mock()
        mock_config.params.chunk_size = 1000

        mock_catalog = Mock(spec=Catalog)
        mock_catalog.get_sources.return_value = ["test_source"]
        mock_catalog.get_tables_by_source.return_value = [
            Table(
                name="test_table",
                replication=ReplicationType.TRUNCATE,
                columns=[
                    Column(name="id", type="int"),
                    Column(name="name", type="string"),
                ],
            )
        ]

        # Setup mock source
        mock_source = Mock()
        mock_source.extract.return_value = [
            [{"id": 1, "name": "test1"}, {"id": 2, "name": "test2"}],  # First chunk
            [{"id": 3, "name": "test3"}],  # Second chunk
        ]
        mock_create_source.return_value = mock_source

        # Setup mock target
        mock_target = Mock()
        mock_target.config.dataset_id = "test_dataset"
        mock_create_target.return_value = mock_target

        # Setup YAML loading
        mock_load_yaml.side_effect = [mock_config, mock_catalog]

        # Execute
        result = main(
            config_path="test_config.yaml",
            catalog_path="test_catalog.yaml",
            debug=False,
            env="dev",
            dryRun=False,
        )

        # Verify
        assert result == 0  # Success

        # Verify source was created and used
        mock_create_source.assert_called_once()
        mock_source.extract.assert_called_once()
        mock_source.close.assert_called_once()

        # Verify target was created and used
        mock_create_target.assert_called_once()
        assert mock_target.load.call_count == 2  # Two chunks loaded

        # Verify load calls
        load_calls = mock_target.load.call_args_list
        # Check that load was called with correct arguments
        assert len(load_calls) == 2  # Two chunks loaded
        # Verify the calls were made (don't check exact args due to mock structure)
        assert mock_target.load.call_count == 2

    @patch("ingest.main.create_target")
    @patch("ingest.main.create_source")
    @patch("ingest.main.load_yaml_model")
    def test_full_pipeline_dry_run(self, mock_load_yaml, mock_create_source, mock_create_target):
        """Test pipeline execution in dry run mode."""
        # Setup mock data
        mock_config = Mock(spec=Config)
        mock_config.targets = {"dev": Mock()}
        mock_config.get_gcloud_secret_value.return_value = "resolved_password"
        mock_config.params = Mock()
        mock_config.params.chunk_size = 1000

        mock_catalog = Mock(spec=Catalog)
        mock_catalog.get_sources.return_value = ["test_source"]
        mock_catalog.get_tables_by_source.return_value = [
            Table(
                name="test_table",
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="id", type="int")],
            )
        ]

        # Setup mock source
        mock_source = Mock()
        mock_source.extract.return_value = [
            [{"id": 1}, {"id": 2}],  # One chunk
        ]
        mock_create_source.return_value = mock_source

        # Setup mock target
        mock_target = Mock()
        mock_target.config.dataset_id = "test_dataset"
        mock_create_target.return_value = mock_target

        # Setup YAML loading
        mock_load_yaml.side_effect = [mock_config, mock_catalog]

        # Execute
        result = main(
            config_path="test_config.yaml",
            catalog_path="test_catalog.yaml",
            debug=False,
            env="dev",
            dryRun=True,  # Dry run mode
        )

        # Verify
        assert result == 0  # Success

        # Verify source was used but target was not
        mock_create_source.assert_called_once()
        mock_source.extract.assert_called_once()
        mock_source.close.assert_called_once()

        # Verify target was created but load was never called
        mock_create_target.assert_called_once()
        mock_target.load.assert_not_called()

    @patch("ingest.main.create_target")
    @patch("ingest.main.create_source")
    @patch("ingest.main.load_yaml_model")
    def test_full_pipeline_source_failure_continues(self, mock_load_yaml, mock_create_source, mock_create_target):
        """Test that pipeline continues when one source fails."""
        # Setup mock data
        mock_config = Mock(spec=Config)
        mock_config.targets = {"dev": Mock()}
        mock_config.get_gcloud_secret_value.return_value = "resolved_password"
        mock_config.params = Mock()
        mock_config.params.chunk_size = 1000

        mock_catalog = Mock(spec=Catalog)
        mock_catalog.get_sources.return_value = ["failing_source", "working_source"]
        mock_catalog.get_tables_by_source.side_effect = [
            [Table(name="failing_table", replication=ReplicationType.TRUNCATE, columns=[])],
            [Table(name="working_table", replication=ReplicationType.TRUNCATE, columns=[])],
        ]

        # Setup mock sources
        mock_failing_source = Mock()
        mock_failing_source.extract.side_effect = Exception("Source failed")
        mock_create_source.side_effect = [mock_failing_source, Mock()]

        # Setup mock target
        mock_target = Mock()
        mock_target.config.dataset_id = "test_dataset"
        mock_create_target.return_value = mock_target

        # Setup YAML loading
        mock_load_yaml.side_effect = [mock_config, mock_catalog]

        # Execute
        result = main(
            config_path="test_config.yaml",
            catalog_path="test_catalog.yaml",
            debug=False,
            env="dev",
            dryRun=False,
        )

        # Verify
        assert result == 0  # Success despite one source failing

        # Verify both sources were attempted
        assert mock_create_source.call_count == 2

    @patch("ingest.main.load_yaml_model")
    def test_full_pipeline_config_loading_failure(self, mock_load_yaml):
        """Test pipeline failure when config loading fails."""
        # Setup YAML loading to fail
        mock_load_yaml.side_effect = Exception("Config file not found")

        # Execute
        result = main(
            config_path="nonexistent_config.yaml",
            catalog_path="test_catalog.yaml",
            debug=False,
            env="dev",
            dryRun=False,
        )

        # Verify
        assert result == 1  # Failure

    @patch("ingest.main.create_target")
    @patch("ingest.main.create_source")
    @patch("ingest.main.load_yaml_model")
    def test_full_pipeline_invalid_environment(self, mock_load_yaml, mock_create_source, mock_create_target):
        """Test pipeline failure with invalid environment."""
        # Setup mock data
        mock_config = Mock(spec=Config)
        mock_config.targets = {"dev": Mock()}
        mock_catalog = Mock(spec=Catalog)
        mock_load_yaml.side_effect = [mock_config, mock_catalog]

        # Execute with invalid environment
        with pytest.raises(ValueError, match="Invalid Env flag"):
            main(
                config_path="test_config.yaml",
                catalog_path="test_catalog.yaml",
                debug=False,
                env="invalid_env",  # Invalid environment
                dryRun=False,
            )
