from typing import Union
from .base import Target
from .bigquery import BigQueryTarget
from ..core.config import Config

def target_factory(config: Config, env: str) -> Target:
    """Create a target instance based on configuration.

    Args:
        config: IngestConfig instance containing target configurations
        env: Environment name (staging/production)

    Returns:
        Appropriate target instance

    Raises:
        ValueError: If target type is not supported or no targets found for environment
    """
    # Get targets for the specified environment
    if env not in config.targets:
        raise ValueError(f"No targets configured for environment: {env}")
    
    targets = config.targets[env]
    if not targets:
        raise ValueError(f"No targets found for environment: {env}")
    
    # For now, use the first target in the environment
    # TODO: Add logic to select specific target or handle multiple targets
    target_config = targets[0]
    
    if target_config.type in {"bigquery", "bq", "google_bigquery"}:
        return BigQueryTarget(target_config, env)
    else:
        raise ValueError(f"Unsupported target type: {target_config.type}")


TargetUnion = Union[BigQueryTarget]

__all__ = [
    "Target",
    "BigQueryTarget",
    "target_factory",
    "TargetUnion",
]
