from typing import Dict, Any, Union
from .base import Target, TargetConfig
from .bigquery import BigQueryTarget


def target_factory(config: Dict[str, Any], env: str) -> Target:
    """Create a target instance based on configuration.

    Args:
        config: Configuration dictionary

    Returns:
        Appropriate target instance

    Raises:
        ValueError: If target type is not supported
    """
    target_type = config.get("type", "").lower()

    if target_type in {"bigquery", "bq", "google_bigquery"}:
        return BigQueryTarget(config, env)
    else:
        raise ValueError(f"Unsupported target type: {target_type}")


TargetUnion = Union[BigQueryTarget]

__all__ = [
    "Target",
    "TargetConfig",
    "BigQueryTarget",
    "target_factory",
    "TargetUnion",
]
