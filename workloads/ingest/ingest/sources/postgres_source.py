from typing import Optional, Any
from .source import Source
from .extractor import BaseExtractor
from ..core.config import SourceConfig
import logging

logger = logging.getLogger(__name__)


class PostgresSource(Source):
    """PostgreSQL data source implementation."""
    
    def __init__(self, config: SourceConfig, env: str):
        """Initialize PostgreSQL source with validated configuration."""
        super().__init__(config=config, env=env)
        
        # Create the extractor internally using base class method
        connection_string = self._build_connection_string(
            dialect="postgresql+psycopg2", 
            default_port=5432
        )
        self.extractor = BaseExtractor(connection_string)
        logger.debug("Initialized PostgreSQL source with BaseExtractor")

    def connect(self) -> None:
        """Connect to the PostgreSQL source."""
        # Delegate to the extractor
        self.extractor.connect()
        logger.debug("PostgreSQL source connected via BaseExtractor")
    
    def extract(self, table: str, limit: Optional[int] = None) -> Any:
        """Extract data from PostgreSQL table."""
        # Delegate to the extractor
        return self.extractor.extract(table, limit)
    
    def validate_connection(self) -> bool:
        """Validate PostgreSQL connection."""
        # Delegate to the extractor
        return self.extractor.validate_connection()