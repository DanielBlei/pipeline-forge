"""Secret handling tests for the ingest workload.

These tests demonstrate:
- Secret resolution and validation
- Production security patterns
- Error handling for secret operations
- Integration with Google Cloud Secret Manager
"""

import pytest
from unittest.mock import Mock, patch

from ingest.helpers.secret_handler import get_gcloud_secret
from ingest.core.config import Config, RuntimeParams, SecretConfig, SecretProvider, BigQueryTargetConfig


class TestSecretHandler:
    """Test secret handler functionality."""

    @patch("ingest.helpers.secret_handler.secretmanager.SecretManagerServiceClient")
    def test_get_gcloud_secret_success(self, mock_client_class):
        """Test successful secret retrieval."""
        # Setup
        mock_client = Mock()
        mock_secret = Mock()
        mock_secret.payload.data.decode.return_value = "secret_value_123"
        mock_client.access_secret_version.return_value = mock_secret
        mock_client_class.return_value = mock_client

        # Execute
        result = get_gcloud_secret("projects/123/secrets/test-secret/versions/latest")

        # Verify
        assert result == "secret_value_123"
        mock_client.access_secret_version.assert_called_once_with(
            name="projects/123/secrets/test-secret/versions/latest"
        )

    @patch("ingest.helpers.secret_handler.secretmanager.SecretManagerServiceClient")
    def test_get_gcloud_secret_client_error(self, mock_client_class):
        """Test secret retrieval with client error."""
        # Setup
        mock_client = Mock()
        mock_client.access_secret_version.side_effect = Exception("Permission denied")
        mock_client_class.return_value = mock_client

        # Execute & Verify
        with pytest.raises(ValueError, match="Failed to get secret"):
            get_gcloud_secret("projects/123/secrets/test-secret/versions/latest")

    @patch("ingest.helpers.secret_handler.secretmanager.SecretManagerServiceClient")
    def test_get_gcloud_secret_empty_value(self, mock_client_class):
        """Test secret retrieval with empty value."""
        # Setup
        mock_client = Mock()
        mock_secret = Mock()
        mock_secret.payload.data.decode.return_value = None
        mock_client.access_secret_version.return_value = mock_secret
        mock_client_class.return_value = mock_client

        # Execute & Verify
        with pytest.raises(ValueError, match="Secret .* is empty"):
            get_gcloud_secret("projects/123/secrets/empty-secret/versions/latest")


class TestConfigSecretIntegration:
    """Test secret integration with configuration."""

    @patch("ingest.core.config.get_gcloud_secret")
    def test_config_secret_resolution_with_secret_path(self, mock_get_secret):
        """Test secret resolution when secret_path is provided."""
        # Setup
        mock_get_secret.return_value = "resolved_password"

        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[
                SecretConfig(
                    name="test-secret",
                    provider=SecretProvider.GOOGLE_SECRET_MANAGER,
                    secret_path="projects/123/secrets/test-secret/versions/latest",
                )
            ],
            sources={},
            targets={},
        )

        # Execute
        result = config.get_gcloud_secret_value("test-secret", "dev")

        # Verify
        assert result == "resolved_password"
        mock_get_secret.assert_called_once_with("projects/123/secrets/test-secret/versions/latest")

    @patch("ingest.core.config.get_gcloud_secret")
    def test_config_secret_resolution_without_secret_path(self, mock_get_secret):
        """Test secret resolution when secret_path is not provided."""
        # Setup
        mock_get_secret.return_value = "resolved_password"

        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[
                SecretConfig(
                    name="test-secret",
                    provider=SecretProvider.GOOGLE_SECRET_MANAGER,
                    # secret_path is None, should use target config
                )
            ],
            sources={},
            targets={
                "dev": BigQueryTargetConfig(
                    name="test_target",
                    project_id="test-project",
                    project_number=123456789,
                    dataset_id="test_dataset",
                )
            },
        )

        # Execute
        result = config.get_gcloud_secret_value("test-secret", "dev")

        # Verify
        assert result == "resolved_password"
        mock_get_secret.assert_called_once_with("projects/123456789/secrets/test-secret/versions/latest")

    def test_config_secret_not_found_in_secrets_list(self):
        """Test error when secret name not found in secrets list."""
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

    def test_config_secret_wrong_provider(self):
        """Test error when secret provider is not Google Secret Manager."""
        # Create a real Config object but manually set a secret with wrong provider
        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[
                SecretConfig(
                    name="test-secret",
                    provider=SecretProvider.GOOGLE_SECRET_MANAGER,  # Start with valid provider
                )
            ],
            sources={},
            targets={},
        )

        # Manually change the provider to simulate wrong provider
        config.secrets[0].provider = "aws"

        # Execute & Verify
        with pytest.raises(ValueError, match="Secret test-secret is not a Gcloud Secret"):
            config.get_gcloud_secret_value("test-secret", "dev")

    def test_config_secret_missing_target_config(self):
        """Test error when target config is missing for secret resolution."""
        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[
                SecretConfig(
                    name="test-secret",
                    provider=SecretProvider.GOOGLE_SECRET_MANAGER,
                    # secret_path is None, but no target config
                )
            ],
            sources={},
            targets={},  # Empty targets
        )

        # Execute & Verify
        with pytest.raises(ValueError, match="Target config not found for environment"):
            config.get_gcloud_secret_value("test-secret", "dev")

    @patch("ingest.core.config.get_gcloud_secret")
    def test_config_secret_resolution_with_custom_version(self, mock_get_secret):
        """Test secret resolution with custom version."""
        # Setup
        mock_get_secret.return_value = "resolved_password"

        config = Config(
            version="1.0.0",
            params=RuntimeParams(),
            secrets=[
                SecretConfig(
                    name="test-secret",
                    provider=SecretProvider.GOOGLE_SECRET_MANAGER,
                    secret_path="projects/123/secrets/test-secret/versions/custom",
                )
            ],
            sources={},
            targets={},
        )

        # Execute
        result = config.get_gcloud_secret_value("test-secret", "dev", version="custom")

        # Verify
        assert result == "resolved_password"
        mock_get_secret.assert_called_once_with("projects/123/secrets/test-secret/versions/custom")
