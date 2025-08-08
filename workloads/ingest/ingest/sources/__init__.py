from typing import Union, Dict, Any
from .base import Source, SourceConfig
from .mysql import MySQLSource
from .postgres import PostgresSource

def source_factory(config: Dict[str, Any], env: str) -> Source:
    """Create a source instance based on configuration.
    
    Args:
        config: Configuration dictionary
        
    Returns:
        Appropriate source instance
        
    Raises:
        ValueError: If source type is not supported
    """
    source_type = config.get("type", "").lower()

    
    if source_type == "mysql":
        return MySQLSource(config, env)
    elif source_type == "postgres":
        return PostgresSource(config, env)
    else:
        raise ValueError(f"Unsupported source type: {source_type}")


# Union type for type hints
SourceUnion = Union[MySQLSource, PostgresSource]


__all__ = [
    "Source",
    "SourceConfig", 
    "MySQLSource",
    "MySQLConfig",
    "PostgresSource", 
    "PostgresConfig",
    "source_factory",
    "SourceUnion"
]

