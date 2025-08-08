from typing import Optional, Any, Dict
from pydantic import BaseModel, Field


# Staging and production targets
class TargetConfig(BaseModel):
    """Base configuration for data targets (destinations)."""
    staging: list[Dict[str, Any]] = Field(default_factory=list, description="List of target configs for staging environment", min_items=1)
    production: list[Dict[str, Any]] = Field(default_factory=list, description="List of target configs for production environment", min_items=1)


class Target(BaseModel):
    """Base target class with Pydantic validation."""
    config: TargetConfig = Field(..., description="Target configuration")

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
