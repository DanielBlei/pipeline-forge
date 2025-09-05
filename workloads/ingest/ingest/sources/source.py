from typing import Optional, Protocol, runtime_checkable, Iterator
from pydantic import BaseModel, Field, ConfigDict
from ..core.config import SourceConfig
from ..core.catalog import Table
from ..extractors import BaseExtractor
import logging

logger = logging.getLogger(__name__)


@runtime_checkable
class SourceInterface(Protocol):
    """Protocol defining the contract for any data source."""

    config: SourceConfig
    extractor: Optional[BaseExtractor]

    def connect(self) -> None: ...
    def extract(self, table: Table, chunk_size: int, limit: Optional[int] = None) -> Iterator[dict]: ...
    def validate_connection(self) -> bool: ...
    def close(self) -> None: ...


class Source(BaseModel):
    """Base implementation of the Source protocol.

    This provides a concrete implementation that can be inherited from
    or used as a reference for implementing the Source protocol.
    """

    config: SourceConfig = Field(..., description="Source configuration")
    extractor: Optional[BaseExtractor] = None
    model_config = ConfigDict(arbitrary_types_allowed=True)

    def connect(self) -> None:
        """Connect to the source."""
        if self.extractor:
            self.extractor.connect()

    def extract(self, table: Table, chunk_size: int, limit: Optional[int] = None) -> Iterator[dict]:
        """Extract data from the source using SQLAlchemy with streaming.

        Args:
            table: Table object to extract from
            chunk_size: Size of data chunks to extract
            limit: Optional limit on number of rows

        Returns:
            Iterator yielding data extracted from the source
        """
        if not self.extractor:
            raise RuntimeError("Extractor not initialized")
        return self.extractor.extract(table=table, chunk_size=chunk_size, limit=limit)

    def validate_connection(self) -> bool:
        """Validate connection to the source.

        Returns:
            True if connection is valid
        """
        if not self.extractor:
            return False
        return self.extractor.validate_connection()

    def close(self) -> None:
        """Close the source connection and clean up resources."""
        if self.extractor:
            self.extractor.close()
