from .target import Target, TargetInterface
from .bigquery import BigQueryTarget
from ..core.config import BigQueryTarget as BigQueryTargetConfig


def create_target(target_config: BigQueryTargetConfig, env: str) -> TargetInterface:
    """Create a target instance directly based on configuration.

    Args:
        target_config: Target configuration object
        env: Environment name (staging/production)

    Returns:
        Appropriate target instance implementing TargetInterface

    Raises:
        ValueError: If target type is not supported
    """
    if target_config.type in {"bigquery", "bq", "google_bigquery"}:
        return BigQueryTarget(target_config, env)
    else:
        raise ValueError(f"Unsupported target type: {target_config.type}")


__all__ = [
    # Target implementations
    "Target",
    "BigQueryTarget",
    # Direct creation function
    "create_target",
    # Protocol interface
    "TargetInterface",
]
