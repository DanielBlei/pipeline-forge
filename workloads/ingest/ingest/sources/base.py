from typing import Optional, Any, Dict
from pydantic import BaseModel, Field

# Staging and production sources
class SourceConfig(BaseModel):
        staging: list[Dict[str, Any]] = Field(default_factory=list, description="List of source configs for staging environment", min_items=1)
        production: list[Dict[str, Any]] = Field(default_factory=list, description="List of source configs for production environment", min_items=1)

class Source(BaseModel):
    """Base source class with Pydantic validation."""
    config: SourceConfig = Field(..., description="Source configuration")
    
    def extract(self, table: str, limit: Optional[int] = None) -> Any:
        """Extract data from the source.
        
        Args:
            table: Name of the table to extract
            limit: Optional limit on number of rows
            
        Returns:
            Extracted data
        """
        raise NotImplementedError("Subclasses must implement extract method")
    
    def validate_connection(self) -> bool:
        """Validate connection to the source.
        
        Returns:
            True if connection is valid
        """
        raise NotImplementedError("Subclasses must implement validate_connection method")
