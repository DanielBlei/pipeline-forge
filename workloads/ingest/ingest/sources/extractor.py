from typing import Optional, Any, Dict, List
from sqlalchemy import create_engine, text, Engine
from sqlalchemy.exc import SQLAlchemyError
import logging

logger = logging.getLogger(__name__)


class BaseExtractor:
    """Base extraction class using SQLAlchemy for database-agnostic data extraction."""
    
    def __init__(self, connection_string: str, **engine_kwargs):
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
            
        except SQLAlchemyError as e:
            logger.error(f"Failed to create engine: {e}")
            raise
    
    def validate_connection(self) -> bool:
        """Validate the database connection.
        
        Returns:
            True if connection is valid, False otherwise
        """
        if not self.engine:
            try:
                self.connect()
            except SQLAlchemyError:
                return False
        
        try:
            with self.engine.connect() as conn:
                # Execute a simple query to test the connection
                result = conn.execute(text("SELECT 1"))
                result.fetchone()
                logger.debug("Database connection validated successfully")
                return True
                
        except SQLAlchemyError as e:
            logger.error(f"Database connection validation failed: {e}")
            return False
    
    def extract(self, table: str, limit: Optional[int] = None) -> List[Dict[str, Any]]:
        """Extract data from a table using SQLAlchemy.
        
        Args:
            table: Name of the table to extract from
            limit: Optional limit on number of rows to extract
            
        Returns:
            List of dictionaries representing the extracted data
            
        Raises:
            SQLAlchemyError: If extraction fails
        """
        if not self.engine:
            self.connect()
        
        try:
            with self.engine.connect() as conn:
                # Build the SELECT query
                query = f"SELECT * FROM {table}"
                if limit:
                    query += f" LIMIT {limit}"
                
                logger.info(f"Executing query: {query}")
                
                # Execute the query
                result = conn.execute(text(query))
                
                # Convert to list of dictionaries
                data = []
                for row in result:
                    data.append(dict(row._mapping))
                
                logger.info(f"Successfully extracted {len(data)} rows from {table}")
                return data
                
        except SQLAlchemyError as e:
            logger.error(f"Failed to extract data from {table}: {e}")
            raise
    
    def close(self) -> None:
        """Close the engine and clean up resources."""
        if self.engine:
            self.engine.dispose()
            logger.info("Engine disposed successfully")
