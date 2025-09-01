from typing import Union
from .source import Source
from .mysql_source import MySQLSource
from .postgres_source import PostgresSource
from ..core.config import Config

def source_factory(config: Config, source_name: str, env: str) -> Source:
    """Create a source instance based on configuration.
    
    Args:
        config: IngestConfig instance containing source configurations
        source_name: Name of the source to create
        env: Environment name (staging/production)
        
    Returns:
        Appropriate source instance
        
    Raises:
        ValueError: If source type is not supported or no sources found for environment
    """
    # Get sources for the specified environment
    if env not in config.sources:
        raise ValueError(f"No sources configured for environment: {env}")
    
    source_config = config.sources.get(env).get(source_name)
    if not source_config:
        available_sources = list(config.sources[env].keys())
        raise ValueError(f"Source '{source_name}' not found in environment '{env}'. Available sources: {available_sources}")
    
    if source_config.type.value == "mysql":
        return MySQLSource(source_config, env)
    elif source_config.type.value == "postgres":
        return PostgresSource(source_config, env)
    else:
        raise ValueError(f"Unsupported source type: {source_config.type}")


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

