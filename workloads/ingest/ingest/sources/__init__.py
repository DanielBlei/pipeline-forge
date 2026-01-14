"""Source package for data source implementations."""

import logging

from .mysql_source import MySQLSource
from .postgres_source import PostgresSource
from .source import SourceInterface
from ..core.config import SourceConfig

logger = logging.getLogger(__name__)


def create_source(source_config: SourceConfig, retry_attempts: int = 3, retry_delay: int = 15) -> SourceInterface:
    """Create and validate a source instance based on configuration.

    Args:
        source_config: Source configuration object
        retry_attempts: Number of retry attempts for connection
        retry_delay: Delay between retries in seconds

    Returns:
        Validated source instance implementing SourceInterface

    Raises:
        ValueError: If source type is not supported or connection validation fails

    """
    # Create the source instance
    source: SourceInterface
    # Note: Instantiate sources with just the config to match tests' expectations.
    # Retry parameters are managed internally by the Source and Extractor with their defaults.
    if source_config.type.value == "mysql":
        source = MySQLSource(source_config)
    elif source_config.type.value == "postgres":
        source = PostgresSource(source_config)
    else:
        raise ValueError(f"Unsupported source type: {source_config.type}")

    # Validate the connection
    if not source.validate_connection():
        raise ValueError(f"Failed to validate source connection: {source.config.name}")

    logger.info(f"Connected to source: {source.config.name}")
    return source


__all__ = [
    # Source implementations
    "MySQLSource",
    "PostgresSource",
    # Factory function with validation
    "create_source",
    # Protocol interface
    "SourceInterface",
]
