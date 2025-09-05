import logging
from typing import Optional


def setup_logging(name: Optional[str] = None, debug: bool = False):
    """Setup logging for the application.

    Args:
        name: Logger name (usually __package__ or __name__)
        debug: Enable debug logging

    Returns:
        Logger instance
    """
    # Get the specific logger if name provided, otherwise root logger
    logger = logging.getLogger(name) if name else logging.getLogger()

    if debug:
        logger.setLevel(logging.DEBUG)
        formatter = logging.Formatter("%(asctime)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s")
    else:
        logger.setLevel(logging.INFO)
        formatter = logging.Formatter("%(asctime)s - %(levelname)s - %(message)s")

    # Only configure if this logger doesn't have handlers
    # This prevents duplicate handlers in the hierarchy
    if not logger.hasHandlers():
        handler = logging.StreamHandler()
        handler.setFormatter(formatter)
        logger.addHandler(handler)

        # Prevent propagation to root logger to avoid duplicates
        if name:
            logger.propagate = False

    return logger
