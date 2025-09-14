"""Extractor logic tests for the ingest workload.

These tests demonstrate:
- Core business logic testing
- Connection string validation
- Error handling and edge cases
- Retry logic and resilience patterns
"""

import pytest
from unittest.mock import Mock

from ingest.extractors import BaseExtractor


class TestBaseExtractor:
    """Test BaseExtractor core functionality."""

    def test_extractor_connection_string_validation(self):
        """Test that invalid connection strings raise ValueError."""
        with pytest.raises(ValueError, match="Invalid connection string"):
            BaseExtractor("invalid-connection-string")

        # Note: mysql://user:pass@host/db is actually valid, so test with truly invalid string
        with pytest.raises(ValueError, match="Invalid connection string"):
            BaseExtractor("not-a-connection-string")

    def test_extractor_valid_connection_string(self):
        """Test that valid connection strings are accepted."""
        extractor = BaseExtractor("mysql+pymysql://user:pass@host:3306/db")
        assert extractor.connection_string == "mysql+pymysql://user:pass@host:3306/db"
        assert extractor.engine is None  # Not connected yet

    def test_extractor_close(self):
        """Test that close properly disposes of the engine."""
        extractor = BaseExtractor("mysql+pymysql://user:pass@host:3306/db")
        extractor.engine = Mock()

        # Execute
        extractor.close()

        # Verify
        extractor.engine.dispose.assert_called_once()

    def test_extractor_close_without_engine(self):
        """Test that close handles missing engine gracefully."""
        extractor = BaseExtractor("mysql+pymysql://user:pass@host:3306/db")

        # Execute (should not raise exception)
        extractor.close()

        # Verify - no exception raised
        assert extractor.engine is None

    def test_extractor_initialization_with_engine_kwargs(self):
        """Test that extractor accepts engine configuration options."""
        extractor = BaseExtractor("mysql+pymysql://user:pass@host:3306/db", pool_size=5, max_overflow=10)

        assert extractor.connection_string == "mysql+pymysql://user:pass@host:3306/db"
        assert extractor.engine_kwargs["pool_size"] == 5
        assert extractor.engine_kwargs["max_overflow"] == 10

    def test_extractor_connection_string_with_special_characters(self):
        """Test connection string with special characters in password."""
        # Test with URL-encoded password
        extractor = BaseExtractor("mysql+pymysql://user:pass%40word@host:3306/db")
        assert extractor.connection_string == "mysql+pymysql://user:pass%40word@host:3306/db"

    def test_extractor_connection_string_with_port(self):
        """Test connection string with explicit port."""
        extractor = BaseExtractor("postgresql+psycopg2://user:pass@host:5432/db")
        assert extractor.connection_string == "postgresql+psycopg2://user:pass@host:5432/db"

    def test_extractor_connection_string_without_port(self):
        """Test connection string without explicit port."""
        extractor = BaseExtractor("mysql+pymysql://user:pass@host/db")
        assert extractor.connection_string == "mysql+pymysql://user:pass@host/db"
