"""Simple hello world test to ensure pytest works."""


def test_hello_world():
    """A basic test that always passes."""
    assert True


def test_simple_math():
    """Another simple test with basic math."""
    assert 1 + 1 == 2


def test_string_operations():
    """Test basic string operations."""
    hello = "Hello"
    world = "World"
    assert f"{hello} {world}" == "Hello World"
