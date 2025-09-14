"""BigQuery target implementation."""

from ..core.catalog import ReplicationType
from .target import Target
from typing import Optional
from ..core.config import BigQueryTargetConfig

import logging
from google.cloud import bigquery
from pydantic import Field

logger = logging.getLogger(__name__)

map_replication_type_to_write_disposition = {
    ReplicationType.TRUNCATE: bigquery.WriteDisposition.WRITE_TRUNCATE,
    ReplicationType.APPEND: bigquery.WriteDisposition.WRITE_APPEND,
    ReplicationType.UPSERT: bigquery.WriteDisposition.WRITE_APPEND,
}


class BigQueryTarget(Target):
    """Google BigQuery target implementation."""

    client: Optional[bigquery.Client] = Field(default=None, exclude=True)

    def __init__(self, config: BigQueryTargetConfig, **data):
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

    def _get_write_disposition(self, write_disposition: ReplicationType) -> str:
        """Get the write disposition for the data."""
        if write_disposition not in map_replication_type_to_write_disposition:
            raise ValueError(f"Invalid write disposition: {write_disposition}")
        return map_replication_type_to_write_disposition[write_disposition]

    def load(self, data: list[dict], target_table: str, write_disposition: ReplicationType) -> None:
        """Load data into BigQuery.

        Args:
            data: The data to be ingested
            target_table: Target table name for the data
            write_disposition: Write disposition for the data

        """
        # Call Target base class validation logic first
        super().load(data, target_table, write_disposition)
        job_config = bigquery.LoadJobConfig(
            write_disposition=self._get_write_disposition(write_disposition),
        )
        if not self.client:
            self.client = self._build_client()

        if not self.client:
            raise RuntimeError("BigQuery client not initialized")

        self.client.load_table_from_json(
            json_rows=data,
            destination=target_table,
            job_config=job_config,
        )

        logger.debug(f"Loading data to BigQuery table: {target_table}")
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
