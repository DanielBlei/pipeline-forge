from typing import Optional, Protocol, runtime_checkable, Any
from pydantic import BaseModel, Field, ConfigDict
from ..core.config import TargetTypes
import logging

logger = logging.getLogger(__name__)


@runtime_checkable
class TargetInterface(Protocol):
    """Protocol defining the contract for any data target."""

    config: TargetTypes

    def load(self, data: Any, table: Optional[str] = None) -> None: ...
    def validate_connection(self) -> bool: ...


class Target(BaseModel):
    """Base implementation of the Target protocol.

    This provides a concrete implementation that can be inherited from
    or used as a reference for implementing the Target protocol.
    """

    config: TargetTypes = Field(..., description="Target configuration")
    model_config = ConfigDict(arbitrary_types_allowed=True)

    def load(self, data: Any, table: Optional[str] = None) -> None:
        """Load data into the target destination.

        Args:
            data: The data to be ingested (format may vary by implementation)
            table: Optional override for the target table or collection name
        """
        raise NotImplementedError("Subclasses must implement load method")

    def validate_connection(self) -> bool:
        """Validate connection to the target destination.

        Returns:
            True if connection is valid
        """
        raise NotImplementedError("Subclasses must implement validate_connection method")
