"""Secret Handler for the ingest workload"""
import logging
from google.cloud import secretmanager

logger = logging.getLogger(__name__)

def get_gcloud_secret(secret_name: str) -> str:
    """Get a secret from the Gcloud Secret Manager"""
    try:
        client = secretmanager.SecretManagerServiceClient()
        secret = client.access_secret_version(name=secret_name)
        secret_data = secret.payload.data.decode("UTF-8")
        if secret_data is None:
            raise ValueError(f"Secret {secret_name} is empty")
        return secret_data
    except Exception as e:
        raise ValueError(f"Failed to get secret {secret_name}: {e}")
