"""Retry decorator for function on exception with a delay between retries."""

from time import sleep
from functools import wraps
import logging

logger = logging.getLogger(__name__)


def retry_on_exception(retries: int = 3, delay: int = 15, retries_attr="retry_attempts", delay_attr="retry_delay"):
    """Universal retry decorator that can use static values or instance attributes.

    This decorator can be used in two ways:
    1. Static retry: @retry_on_exception(retries=5, delay=10)
    2. Instance-based retry: @retry_on_exception() - reads from instance attributes

    Args:
        retries: Static retry count (used if instance doesn't have retries_attr)
        delay: Static delay in seconds (used if instance doesn't have delay_attr)
        retries_attr: Name of instance attribute for retry count (default: 'retry_attempts')
        delay_attr: Name of instance attribute for delay (default: 'retry_delay')
    """

    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            # Determine if this is a method (has 'self') or a function
            if args and hasattr(args[0], retries_attr):
                # It's a method with instance attributes - use them
                self = args[0]
                retries_val = getattr(self, retries_attr, retries)
                delay_val = getattr(self, delay_attr, delay)
            else:
                # It's a function or method without instance attributes - use static values
                retries_val = retries
                delay_val = delay

            for attempt in range(retries_val):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    attempt_num = attempt + 1
                    logger.error(f"Attempt {attempt_num}/{retries_val} failed: {e}")
                    if attempt < retries_val - 1:
                        logger.info(f"Retrying in {delay_val} seconds...")
                        sleep(delay_val)
                    else:
                        logger.error(f"All {retries_val} attempts failed. Final attempt:")
            # Final attempt, let exception propagate
            return func(*args, **kwargs)

        return wrapper

    return decorator
