from typing import Optional, Any
from rich.table import Table
from sqlalchemy import create_engine, text, Engine
from sqlalchemy.exc import SQLAlchemyError
import logging

from ..core.config import Config
from ..helpers.retry import retry_on_exception

logger = logging.getLogger(__name__)


class BaseExtractor:
    """Base extraction class using SQLAlchemy for database-agnostic data extraction."""

    def __init__(self, connection_string: str, config: Optional[Config] = None, **engine_kwargs):
        """Initialize the extractor with a database connection string.
        
        Args:
            connection_string: SQLAlchemy connection string
            **engine_kwargs: Additional engine configuration options
        """
        if "://" not in connection_string:
            raise ValueError("Invalid connection string: must be a valid SQLAlchemy URL like 'dialect+driver://user:pass@host[:port]/db'")
        
        self.connection_string = connection_string
        self.engine: Optional[Engine] = None
        self.engine_kwargs = engine_kwargs
    
    @retry_on_exception(retries=5, delay=30)
    def connect(self) -> None:
        """Connect to the database and create the engine."""
        try:
            # Set default engine options for better performance and reliability
            default_options = {
                "pool_size": 10,
                "max_overflow": 20,
                "pool_recycle": 3600,  # Recycle connections after 1 hour
                "pool_pre_ping": True,  # Validate connections before use
                "echo": False,  # Set to True for debugging
            }
            
            # Merge with user-provided options
            engine_options = {**default_options, **self.engine_kwargs}
            
            logger.info("Connecting to database...")
            self.engine = create_engine(self.connection_string, **engine_options)
            logger.debug("Engine created successfully")
            
        except Exception as e:
            logger.error(f"Failed to create engine: {e}")
            raise
    
    @retry_on_exception(retries=5, delay=30)
    def validate_connection(self) -> bool:
        """Validate the database connection.
        
        Returns:
            True if connection is valid, False otherwise
        """
        if not self.engine:
            self.connect()  # This will retry if it fails
        
        try:
            with self.engine.connect() as conn:
                # Execute a simple query to test the connection
                result = conn.execute(text("SELECT 1"))
                result.fetchone()
                logger.debug("Database connection validated successfully")
                return True
                
        except SQLAlchemyError as e:
            logger.error(f"Database connection validation failed: {e}")
            raise
    
    @retry_on_exception()
    def extract(self, table: Table, chunk_size: int = 1000, limit: Optional[int] = None) -> Any:
        """Extract data from a table using SQLAlchemy with streaming for memory efficiency.
        
        Args:
            table: Name of the table to extract from
            limit: Optional limit on number of rows to extract
            
        Yields:
            Dict[str, Any] representing each row from the table
            
        Raises:
            SQLAlchemyError: If extraction fails
        """
        if not self.engine:
            self.connect()
        
        try:
            with self.engine.connect() as conn:
                # Build the SELECT query
                columns = ",".join([column.name for column in table.columns])
                query = f"SELECT {columns} FROM {table.name}"
                if limit:
                    query += f" LIMIT {limit}"
                
                logger.info(f"Executing streaming query: {query}")
                
                # Execute with cursor for memory efficiency
                result = conn.execution_options(yield_per=chunk_size).execute(text(query))
                
                # Yield rows one by one for memory efficiency
                for row in result:
                    yield dict(row._mapping)
                
                logger.info(f"Successfully completed streaming extraction from {table.name}")
                
        except SQLAlchemyError as e:
            logger.error(f"Failed to extract data from {table.name}: {e}")
            raise   
    
    def close(self) -> None:
        """Close the engine and clean up resources."""
        if self.engine:
            self.engine.dispose()
            logger.info("Engine disposed successfully")
