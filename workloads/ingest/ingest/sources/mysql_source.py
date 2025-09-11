from typing import Optional, Iterator
from .source import Source
from ..extractors import BaseExtractor
from ..core.config import SourceConfig
from ..core.catalog import Table
import logging

logger = logging.getLogger(__name__)


class MySQLSource(Source):
    """MySQL data source implementation."""

    def __init__(self, config: SourceConfig):
        """Initialize MySQL source with validated configuration."""
        super().__init__(config=config)
        logger.debug("Initialized MySQL source")

        # Initialize the extractor with the connection string
        connection_string = self.config.build_connection_string(dialect="mysql+pymysql", ext_password=self.config.password, default_port=3306)
        self.extractor = BaseExtractor(connection_string)
        logger.debug("Initialized BaseExtractor")

    def connect(self) -> None:
        """Connect to the MySQL source."""
        super().connect()
        logger.debug("MySQL source connected via BaseExtractor")

    def extract(self, table: Table, chunk_size: int, limit: Optional[int] = None) -> Iterator[dict]:
        """Extract data from MySQL table."""
        if not self.extractor:
            raise RuntimeError("MySQL extractor not initialized")
        return self.extractor.extract(table=table, chunk_size=chunk_size, limit=limit)

    def validate_connection(self) -> bool:
        """Validate MySQL connection."""
        if not self.extractor:
            return False
        return self.extractor.validate_connection()
