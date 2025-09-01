from typing import Optional, Any
from .source import Source
from ..extractors import BaseExtractor
from ..core.config import SourceConfig
import logging

logger = logging.getLogger(__name__)


class MySQLSource(Source):
    """MySQL data source implementation."""
    
    def __init__(self, config: SourceConfig, env: str):
        """Initialize MySQL source with validated configuration."""
        super().__init__(config=config, env=env)
        
        # Create the extractor internally using base class method
        connection_string = self._build_connection_string(
            dialect="mysql+pymysql", 
            default_port=3306
        )
        self.extractor = BaseExtractor(connection_string)   
        logger.debug("Initialized MySQL source with BaseExtractor")

    def connect(self) -> None:
        """Connect to the MySQL source."""
        # Delegate to the extractor
        self.extractor.connect()
        logger.debug("MySQL source connected via BaseExtractor")
    
    def extract(self, table: str, chunk_size: int, limit: Optional[int] = None) -> Any:
        """Extract data from MySQL table."""
        # Delegate to the extractor
        return self.extractor.extract(table=table, chunk_size=chunk_size, limit=limit)
    
    def validate_connection(self) -> bool:
        """Validate MySQL connection."""
        # Delegate to the extractor
        return self.extractor.validate_connection()