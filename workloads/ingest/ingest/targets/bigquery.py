from .target import Target
from typing import Any, Optional
from ..core.config import BigQueryTarget
import logging
from google.cloud import bigquery
from pydantic import Field

logger = logging.getLogger(__name__)


class BigQueryTarget(Target):
    """Google BigQuery target implementation."""

    client: Optional[bigquery.Client] = Field(default=None, exclude=True)

    def __init__(self, config: BigQueryTarget, **data):
        """Initialize BigQuery target with validated configuration."""
        super().__init__(config=config)
        self.client: Optional[bigquery.Client] = None
        logger.debug(f"Initialized BigQuery target for project: {self.config.project_id}")

    def _build_client(self) -> bigquery.Client:
        """Build a BigQuery client."""
        try:
            client = bigquery.Client(project=self.config.project_id)
        except Exception as e:
            logger.error(f"Failed to build BigQuery client: {e}")
            raise e
        return client

    def load(self, data: Any, table: Optional[str] = None) -> None:
        """Load data into BigQuery.

        Args:
            data: The data to be ingested
            table: Optional override for the target table name
        """
        # TODO: Implement loading logic to BigQuery
        logger.debug(f"Loading data to BigQuery table: {table or 'default'}")
        return None

    def validate_connection(self) -> bool:
        """Validate connection to BigQuery.

        Returns:
            True if connection is valid
        """
        if not self.client:
            self.client = self._build_client()

        # check if client is built
        if not self.client:
            return False

        try:
            self.client.query("SELECT 1")
        except Exception as e:
            logger.error(f"Failed to validate BigQuery connection: {e}")
            return False
        logger.debug("Validating BigQuery connection")
        return True
