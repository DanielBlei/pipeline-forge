from typing import Optional, Iterator
from .source import Source
from ..extractors import BaseExtractor
from ..core.config import SourceConfig
from ..core.catalog import Table
import logging

logger = logging.getLogger(__name__)


class PostgresSource(Source):
    """PostgreSQL data source implementation."""

    def __init__(self, config: SourceConfig, retry_attempts: int = 3, retry_delay: int = 15):
        """Initialize PostgreSQL source with validated configuration."""
        super().__init__(config=config, retry_attempts=retry_attempts, retry_delay=retry_delay)
        logger.debug("Initialized PostgreSQL source")

        # Initialize the extractor with the connection string
        connection_string = self.config.build_connection_string(dialect="postgresql+psycopg2", default_port=5432)
        self.extractor = BaseExtractor(connection_string, retry_attempts=retry_attempts, retry_delay=retry_delay)
        logger.debug("Initialized BaseExtractor")

    def connect(self) -> None:
        """Connect to PostgreSQL source."""
        super().connect()
        logger.debug("PostgreSQL source connected via BaseExtractor")

    def extract(self, table: Table, chunk_size: int, limit: Optional[int] = None) -> Iterator[dict]:
        """Extract data from PostgreSQL table."""
        if not self.extractor:
            raise RuntimeError("PostgreSQL extractor not initialized")
        return self.extractor.extract(table=table, chunk_size=chunk_size, limit=limit)

    def validate_connection(self) -> bool:
        """Validate PostgreSQL connection."""
        if not self.extractor:
            return False
        return self.extractor.validate_connection()

    def close(self) -> None:
        """Close the source connection and clean up resources."""
        super().close()
        logger.debug("Closed PostgreSQL source connection")
