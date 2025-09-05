import logging

from .target import Target, TargetInterface
from .bigquery import BigQueryTarget
from ..core.config import BigQueryTarget as BigQueryTargetConfig

logger = logging.getLogger(__name__)


def create_target(target_config: BigQueryTargetConfig) -> TargetInterface:
    """Create and validate a target instance based on configuration.

    Args:
        target_config: Target configuration object

    Returns:
        Validated target instance implementing TargetInterface

    Raises:
        ValueError: If target type is not supported or connection validation fails
    """
    # Create the target instance
    if target_config.type in {"bigquery", "bq", "google_bigquery"}:
        target = BigQueryTarget(target_config)
    else:
        raise ValueError(f"Unsupported target type: {target_config.type}")

    # Validate the connection
    if not target.validate_connection():
        raise ValueError(f"Failed to validate target connection: {target.config.name}")

    logger.debug(f"Target is initialized: {target.config.name}")
    return target


__all__ = [
    # Target implementations
    "Target",
    "BigQueryTarget",
    # Factory function with validation
    "create_target",
    # Protocol interface
    "TargetInterface",
]
