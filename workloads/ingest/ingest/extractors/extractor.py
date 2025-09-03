from typing import Optional, Iterator, List
from sqlalchemy import create_engine, text, Engine
from sqlalchemy.exc import SQLAlchemyError
import logging

from ..core.config import Config
from ..core.catalog import Table
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
            raise ValueError(
                "Invalid connection string: must be a valid SQLAlchemy URL like 'dialect+driver://user:pass@host[:port]/db'"
            )

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
    def extract(self, table: Table, chunk_size: int = 1000, limit: Optional[int] = None) -> Iterator[List[dict]]:
        """Extract data from a table using SQLAlchemy with streaming for memory efficiency.

        Args:
            table: Name of the table to extract from
            chunk_size: Number of rows to yield in each chunk
            limit: Optional limit on number of rows to extract

        Yields:
            List[Dict[str, Any]] representing chunks of rows from the table

        Raises:
            SQLAlchemyError: If extraction fails
        """
        if not self.engine:
            self.connect()

        try:
            with self.engine.connect() as conn:
                # Add execution options for guaranteed streaming support with MySQL and Postgres
                conn.execution_options(yield_per=chunk_size, stream_results=True)

                # Build the SELECT query
                columns = ",".join([column.name for column in table.columns])
                query = f"SELECT {columns} FROM {table.name}"
                if limit:
                    query += f" LIMIT {limit}"

                logger.debug(f"Executing streaming query: {query}")
                result = conn.execute(text(query))

                chunk = []
                for row in result:
                    chunk.append(dict(row._mapping))
                    if len(chunk) >= chunk_size:
                        yield chunk
                        chunk = []
                # Yield any remaining rows in the final chunk
                if chunk:
                    yield chunk

        except SQLAlchemyError as e:
            logger.error(f"Failed to extract data from {table.name}: {e}")
            raise

    def close(self) -> None:
        """Close the engine and clean up resources."""
        if self.engine:
            self.engine.dispose()
            logger.debug("Engine disposed successfully")
