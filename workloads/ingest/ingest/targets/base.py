from typing import Optional, Any
from pydantic import BaseModel, Field
from ..core.config import BigQueryTarget as BigQueryTargetConfig


class Target(BaseModel):
    """Base target class with Pydantic validation."""
    config: BigQueryTargetConfig = Field(..., description="Target configuration")
    env: str = Field(..., description="Environment name (staging/production)")

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
