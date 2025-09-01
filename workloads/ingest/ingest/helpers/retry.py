"""Retry decorator for function on exception with a delay between retries."""

from time import sleep
from functools import wraps
import logging

logger = logging.getLogger(__name__)


def retry_on_exception(retries: int = 3, delay: int = 15):
    """Retry function on exception with a delay between retries."""
    def decorator(func):    
        @wraps(func)
        def wrapper(*args, **kwargs):
            for attempt in range(retries):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    attempt_num = attempt + 1
                    logger.error(f"Attempt {attempt_num}/{retries} failed: {e}")
                    if attempt < retries - 1:
                        logger.info(f"Retrying in {delay} seconds...")
                        sleep(delay)
                    else:
                        logger.error(f"All {retries} attempts failed. Final attempt:")
            # Final attempt, let exception propagate
            return func(*args, **kwargs)
        return wrapper
    return decorator