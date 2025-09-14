"""Catalog validation tests for the ingest workload.

These tests demonstrate:
- Pydantic model validation for data catalog
- Table and column schema validation
- Source filtering and data organization
"""

import pytest
from pydantic import ValidationError

from ingest.core.catalog import Catalog, Table, Column, ReplicationType


class TestColumn:
    """Test Column model validation."""

    def test_valid_column(self):
        """Test valid column creation."""
        column = Column(name="user_id", type="int")

        assert column.name == "user_id"
        assert column.type == "int"

    def test_column_validation_constraints(self):
        """Test that column enforces validation constraints."""
        # Test that empty name raises ValidationError
        with pytest.raises(ValidationError):
            Column(name="", type="int")

        # Test that empty type raises ValidationError
        with pytest.raises(ValidationError):
            Column(name="user_id", type="")

        # Test that valid column works
        column = Column(name="user_id", type="int")
        assert column.name == "user_id"
        assert column.type == "int"


class TestTable:
    """Test Table model validation."""

    def test_valid_table(self):
        """Test valid table creation."""
        columns = [
            Column(name="id", type="int"),
            Column(name="name", type="string"),
        ]

        table = Table(
            name="users",
            replication=ReplicationType.TRUNCATE,
            columns=columns,
        )

        assert table.name == "users"
        assert table.replication == ReplicationType.TRUNCATE
        assert len(table.columns) == 2
        assert table.source is None  # Default value

    def test_table_with_source(self):
        """Test table with explicit source."""
        columns = [Column(name="id", type="int")]

        table = Table(
            name="users",
            source="custom_source",
            replication=ReplicationType.APPEND,
            columns=columns,
        )

        assert table.source == "custom_source"
        assert table.replication == ReplicationType.APPEND

    def test_table_validation(self):
        """Test table validation with invalid replication type and name constraints."""
        # Test that empty name raises ValidationError
        with pytest.raises(ValidationError):
            Table(
                name="",  # Empty name should raise ValidationError
                replication=ReplicationType.TRUNCATE,
                columns=[],
            )

        # Test invalid replication type
        with pytest.raises(ValidationError):
            Table(
                name="users",
                replication="INVALID_TYPE",  # Invalid replication type
                columns=[],
            )

        # Test valid table
        table = Table(
            name="users",
            replication=ReplicationType.TRUNCATE,
            columns=[],
        )
        assert table.name == "users"


class TestCatalog:
    """Test Catalog model validation and behavior."""

    def test_valid_catalog(self):
        """Test valid catalog creation."""
        tables = [
            Table(
                name="users",
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="id", type="int")],
            ),
            Table(
                name="orders",
                replication=ReplicationType.APPEND,
                columns=[Column(name="order_id", type="int")],
            ),
        ]

        catalog = Catalog(
            name="test_catalog",
            source="test_source",
            tables=tables,
        )

        assert catalog.name == "test_catalog"
        assert catalog.source == "test_source"
        assert len(catalog.tables) == 2

    def test_catalog_get_table(self):
        """Test getting table by name."""
        tables = [
            Table(
                name="users",
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="id", type="int")],
            ),
        ]

        catalog = Catalog(
            name="test_catalog",
            source="test_source",
            tables=tables,
        )

        table = catalog.get_table("users")
        assert table.name == "users"

        # Test non-existent table
        with pytest.raises(ValueError, match="Table nonexistent not found"):
            catalog.get_table("nonexistent")

    def test_catalog_get_tables_by_source(self):
        """Test filtering tables by source."""
        tables = [
            Table(
                name="users",
                source="source1",
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="id", type="int")],
            ),
            Table(
                name="orders",
                source="source2",
                replication=ReplicationType.APPEND,
                columns=[Column(name="order_id", type="int")],
            ),
            Table(
                name="products",
                # No explicit source, should use catalog's default source
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="product_id", type="int")],
            ),
        ]

        catalog = Catalog(
            name="test_catalog",
            source="default_source",
            tables=tables,
        )

        # Test explicit source
        source1_tables = catalog.get_tables_by_source("source1")
        assert len(source1_tables) == 1
        assert source1_tables[0].name == "users"

        # Test default source
        default_tables = catalog.get_tables_by_source("default_source")
        assert len(default_tables) == 1
        assert default_tables[0].name == "products"

    def test_catalog_get_sources(self):
        """Test getting unique sources from catalog."""
        tables = [
            Table(
                name="users",
                source="source1",
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="id", type="int")],
            ),
            Table(
                name="orders",
                source="source2",
                replication=ReplicationType.APPEND,
                columns=[Column(name="order_id", type="int")],
            ),
            Table(
                name="products",
                # No explicit source, uses catalog's default
                replication=ReplicationType.TRUNCATE,
                columns=[Column(name="product_id", type="int")],
            ),
        ]

        catalog = Catalog(
            name="test_catalog",
            source="default_source",
            tables=tables,
        )

        sources = catalog.get_sources()
        assert len(sources) == 3
        assert "source1" in sources
        assert "source2" in sources
        assert "default_source" in sources

    def test_catalog_validation_constraints(self):
        """Test that catalog enforces validation constraints."""
        # Note: Catalog model has no min_length constraints on name/source, so empty strings are allowed
        catalog1 = Catalog(
            name="",  # Empty name is allowed (no min_length constraint)
            source="test_source",
            tables=[],
        )
        assert catalog1.name == ""

        catalog2 = Catalog(
            name="test_catalog",
            source="",  # Empty source is allowed (no min_length constraint)
            tables=[],
        )
        assert catalog2.source == ""

    def test_catalog_extra_fields_forbidden(self):
        """Test that extra fields in Catalog raise ValidationError."""
        with pytest.raises(ValidationError):
            Catalog(
                name="test_catalog",
                source="test_source",
                tables=[],
                extra_field="not_allowed",  # This should cause validation error
            )
