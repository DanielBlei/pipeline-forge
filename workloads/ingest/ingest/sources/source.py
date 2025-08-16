from typing import Optional, Any
from pydantic import BaseModel, Field, ConfigDict
from ..core.config import SourceConfig
from .extractor import BaseExtractor
import logging

logger = logging.getLogger(__name__)

class Source(BaseModel):
    """Base source class with Pydantic validation and SQLAlchemy extraction."""
    config: SourceConfig = Field(..., description="Source configuration")
    env: str = Field(..., description="Environment name (staging/production)")
    extractor: Optional[BaseExtractor] = None
    model_config = ConfigDict(arbitrary_types_allowed=True)
   
    def _build_connection_string(self, dialect: str, default_port: int) -> str:
        """Build database connection string from config.
        
        Args:
            dialect: SQLAlchemy dialect (e.g., 'mysql+pymysql', 'postgresql+psycopg2')
            default_port: Default port for the database
            
        Returns:
            SQLAlchemy connection string
        """
        host = self.config.connection.host
        port = self.config.connection.port
        database = self.config.connection.database
        username = self.config.connection.username
        password = self.config.connection.password
        
        return f"{dialect}://{username}:{password}@{host}:{port}/{database}"

    def connect(self) -> None:
        """Connect to the source."""
        if self.extractor:
            self.extractor.connect()
        logger.info(f"Connected to source: {self.config.database}")
    
    def extract(self, table: str, limit: Optional[int] = None) -> Any:
        """Extract data from the source using SQLAlchemy.
        
        Args:
            table: Name of the table to extract
            limit: Optional limit on number of rows
            
        Returns:
            Extracted data as list of dictionaries
        """
        if not self.extractor:
            raise RuntimeError("Extractor not initialized")
        return self.extractor.extract(table, limit)
    
    def validate_connection(self) -> bool:
        """Validate connection to the source.
        
        Returns:
            True if connection is valid
        """
        if not self.extractor:
            return False
        return self.extractor.validate_connection()
    
    def close(self) -> None:
        """Close the source connection and clean up resources."""
        if self.extractor:
            self.extractor.close()
        logger.info("Closed source connection")
