"""Catalog class for storing the schema of the data"""

from pydantic import BaseModel, Field, ConfigDict
from enum import Enum


class Column(BaseModel):
    """Column class for storing the schema of the data"""

    name: str
    type: str


class ReplicationType(str, Enum):
    "Replication Method, e.g TRUNCATE, APPEND, UPSERT"

    TRUNCATE = "TRUNCATE"
    APPEND = "APPEND"
    UPSERT = "UPSERT"


class Table(BaseModel):
    """Table class for storing the schema of the data"""

    name: str
    source: str | None = Field(
        default=None,
        description="Source name for this table, matching the source name in the config. If not set, defaults to the catalog's source.",
    )
    replication: ReplicationType
    columns: list[Column]


class Catalog(BaseModel):
    """Catalog class for storing the schema of the data"""

    name: str = Field(description="Catalog name")
    source: str = Field(description="Source name for this catalog, matching the source name in the config")
    tables: list[Table] = Field(description="Tables from various sources")

    model_config = ConfigDict(extra="forbid", validate_assignment=True, str_strip_whitespace=True)

    def get_table(self, table_name: str) -> Table:
        """Get a table by name"""
        for table in self.tables:
            if table.name == table_name:
                return table
        raise ValueError(f"Table {table_name} not found in catalog")

    def get_tables_by_source(self, source_name: str) -> list[Table]:
        """Get all tables from a specific source"""
        return [
            table for table in self.tables if (table.source if table.source is not None else self.source) == source_name
        ]

    def get_sources(self) -> list[str]:
        """Get list of all unique sources in the catalog"""
        return list(set(table.source if table.source is not None else self.source for table in self.tables))
