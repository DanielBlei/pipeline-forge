from typing import Union
from .source import Source
from .mysql_source import MySQLSource
from .postgres_source import PostgresSource
from ..core.config import Config

def source_factory(config: Config, env: str) -> Source:
    """Create a source instance based on configuration.
    
    Args:
        config: IngestConfig instance containing source configurations
        env: Environment name (staging/production)
        
    Returns:
        Appropriate source instance
        
    Raises:
        ValueError: If source type is not supported or no sources found for environment
    """
    # Get sources for the specified environment
    if env not in config.sources:
        raise ValueError(f"No sources configured for environment: {env}")
    
    sources = config.sources[env]
    if not sources:
        raise ValueError(f"No sources found for environment: {env}")
    
    for source in sources:
        if source.type.value == "mysql":
            return MySQLSource(source, env)
        elif source.type.value == "postgres":
            return PostgresSource(source, env)
        else:
            raise ValueError(f"Unsupported source type: {source.type}")


# Union type for type hints
SourceUnion = Union[MySQLSource, PostgresSource]


__all__ = [
    # Source classes
    "Source",
    "MySQLSource",
    "PostgresSource", 

    # Factory function
    "source_factory",

    # Supported source types
    "SourceUnion"
]

