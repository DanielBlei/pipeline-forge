from .mysql_source import MySQLSource
from .postgres_source import PostgresSource
from .source import SourceInterface
from ..core.config import SourceConfig


def create_source(source_config: SourceConfig, env: str) -> SourceInterface:
    """Create a source instance directly based on configuration.

    Args:
        source_config: Source configuration object
        env: Environment name (staging/production)

    Returns:
        Appropriate source instance implementing SourceInterface

    Raises:
        ValueError: If source type is not supported
    """
    if source_config.type.value == "mysql":
        return MySQLSource(source_config, env)
    elif source_config.type.value == "postgres":
        return PostgresSource(source_config, env)
    else:
        raise ValueError(f"Unsupported source type: {source_config.type}")


__all__ = [
    # Source implementations
    "MySQLSource",
    "PostgresSource",
    # Direct creation function
    "create_source",
    # Protocol interface
    "SourceInterface",
]
