import logging


def setup_logging(debug: bool = False):
    logger = logging.getLogger()
    if debug:
        logger.setLevel(logging.DEBUG)
        formatter = logging.Formatter("%(asctime)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s")
    else:
        logger.setLevel(logging.INFO)
        formatter = logging.Formatter("%(asctime)s - %(levelname)s - %(message)s")

    handler = logging.StreamHandler()
    handler.setFormatter(formatter)
    if not logger.hasHandlers():
        logger.addHandler(handler)
    else:
        # Replace existing handlers' formatters
        for h in logger.handlers:
            h.setFormatter(formatter)
    return logger
