"""Base implementation of the Target protocol."""

from typing import Protocol, runtime_checkable
from pydantic import BaseModel, Field, ConfigDict
from ..core.config import TargetTypes
from ..core.catalog import ReplicationType
import logging

logger = logging.getLogger(__name__)


@runtime_checkable
class TargetInterface(Protocol):
    """Protocol defining the contract for any data target."""

    config: TargetTypes

    def load(self, data: list[dict], target_table: str, write_disposition: ReplicationType) -> None:
        """Load data into the target."""
        ...

    def validate_connection(self) -> bool:
        """Validate the connection to the target."""
        ...


class Target(BaseModel):
    """Base implementation of the Target protocol.

    This provides a concrete implementation that can be inherited from
    or used as a reference for implementing the Target protocol.
    """

    config: TargetTypes = Field(..., description="Target configuration")
    model_config = ConfigDict(arbitrary_types_allowed=True)

    def load(self, data: list[dict], target_table: str, write_disposition: ReplicationType) -> None:
        """Load data into the target destination.

        Args:
            data: The data to be ingested
            target_table: Target table or collection name
            write_disposition: How to handle existing data (TRUNCATE, APPEND, UPSERT)

        Raises:
            ValueError: If data is empty or target_table is not provided
            NotImplementedError: If called on base class directly

        """
        if not data:
            raise ValueError("No data to load into target")
        if not target_table or not target_table.strip():
            raise ValueError("Target table name is required")

        if "." not in target_table or len(target_table.split(".")) != 2:
            raise ValueError("Target table name must be in the format 'dataset.table'")

        if not write_disposition:
            raise ValueError("Write disposition is required")

    def validate_connection(self) -> bool:
        """Validate connection to the target destination.

        Returns:
            True if connection is valid

        Note:
            This is a base implementation that should be overridden by subclasses.
            Each target type has different connection validation requirements.

        """
        raise NotImplementedError("Subclasses must implement validate_connection method")
