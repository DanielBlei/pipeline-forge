from pathlib import Path
import sys
from typing import Any

import typer  # type: ignore

from ruamel.yaml import YAML  # type: ignore
from .log import setup_logging
from .sources import create_source
from .core import Config, Catalog
from .targets import create_target, Target

app = typer.Typer()

# Module-level constants for default values
CONFIG_PATH = typer.Option(..., "--config", "-c", help="Path to config file")
CATALOG_PATH = typer.Option(..., "--catalog", "-cat", help="Path to catalog file")
DEBUG_FLAG = typer.Option(False, "--debug", "-d", help="Enable debug mode")
ENV_FLAG = typer.Option("dev", "--env", "-e", help="Environment to use")

logger = setup_logging(debug=DEBUG_FLAG)


@app.command()
def ingest(
    configPath: Path = CONFIG_PATH,
    catalogPath: Path = CATALOG_PATH,
    debug: bool = DEBUG_FLAG,
    env: str = ENV_FLAG,
) -> int:
    """Ingest data from a source database to a target system.

    Args:
        configPath: Path to config file
        catalogPath: Path to catalog file
        debug: Enable debug mode
        env: Environment to use

    Returns:
        int: Exit code (0 for success, 1 for failure)
    """
    try:
        return main(configPath, catalogPath, debug, env)
    except Exception as e:
        logger.error(f"Unexpected error in main: {e}")
        return 1


def main(config_path: Path, catalog_path: Path, debug: bool, env: str) -> int:
    """Main entry point for the ingestion process.

    Args:
        config_path: Path to config file
        catalog_path: Path to catalog file
        debug: Enable debug mode
        env: Environment to use

    Returns:
        int: Exit code (0 for success, 1 for failure)
    """
    # Initialize logger with the debug parameter
    logger = setup_logging(debug=debug)
    logger.info("Starting ingestion process")

    if env.lower() not in ["dev", "prod"]:
        raise ValueError(f"Invalid Env flag: {env} expected `dev` or `prod`")

    try:
        config = load_yaml_model(config_path, Config)
        catalog = load_yaml_model(catalog_path, Catalog)

        target = create_and_validate_target(config.targets[env], env, logger)

        # Process each unique source
        sources = catalog.get_sources()
        for source_name in sources:
            try:
                process_source(source_name, target, config, catalog, env, logger)
            except Exception as e:
                logger.error(f"Failed to process source {source_name}: {e}")
                continue  # TODO: add a flag to halt the process if desired
            logger.info(f"Ingested data from source '{source_name}'")

        logger.info("Ingestion job successfully completed")
        return 0

    except ValueError as e:
        logger.error(f"Configuration error: {e}")
        return 1
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return 1


def process_source(source_name: str, target: Target, config: Config, catalog: Catalog, env: str, logger) -> None:
    """Process one source and all its tables.

    Args:
        source_name: Name of the source to process
        target: Target instance for data loading
        config: Configuration object
        catalog: Catalog object containing table definitions
        env: Environment to use
        logger: Logger instance
    """
    # Get all tables for this source
    tables = catalog.get_tables_by_source(source_name)
    logger.info(f"Processing source '{source_name}' with {len(tables)} tables")

    source_config = config.get_source_config(env, source_name)
    source = create_and_validate_source(source_name, source_config, env, logger)
    try:
        # Process each table from this source
        for table in tables:
            try:
                process_table(source, target, table, config.params.chunk_size, logger)
            except Exception as e:
                logger.error(f"Failed to process table {table.name}: {e}")
                continue  # TODO: add a flag to halt the process if desired
    finally:
        source.close()


def process_table(source, target, table, chunk_size: int, logger) -> None:
    """Process one table's extraction and loading.

    Args:
        source: Source instance for data extraction
        target: Target instance for data loading
        table: Table definition to process
        chunk_size: Size of data chunks to process
        logger: Logger instance
    """
    chunk_count = 0
    for chunk in source.extract(table=table, chunk_size=chunk_size):
        chunk_count += 1
        target.load(chunk, table.name)
        logger.info(f"Loaded chunk {chunk_count} from table '{table.name}'")

    logger.info(f"Completed table '{table.name}': {chunk_count} chunks processed")


def create_and_validate_source(source_name: str, source_config, env: str, logger):
    """Create and validate a source connection.

    Args:
        source_name: Name of the source
        source_config: Source configuration
        env: Environment to use
        logger: Logger instance

    Returns:
        Source instance with validated connection

    Raises:
        ValueError: If source connection validation fails
    """
    source = create_source(source_config, env)
    if not source.validate_connection():
        raise ValueError(f"Failed to validate source connection: {source.config.name}")
    logger.info(f"Connected to source: {source.config.name}")
    return source


def create_and_validate_target(target_config, env: str, logger):
    """Create and validate a target connection.

    Args:
        target_config: Target configuration
        env: Environment to use
        logger: Logger instance

    Returns:
        Target instance with validated connection

    Raises:
        ValueError: If target connection validation fails
    """
    target = create_target(target_config, env)
    if not target.validate_connection():
        raise ValueError(f"Failed to validate target connection: {target.config.name}")
    logger.debug(f"Targer is initialized: {target.config.name}")
    return target


def load_yaml_model(path: Path, model_cls) -> Any:
    """Load and validate a YAML file into a model instance.

    Args:
        path: Path to the YAML file
        model_cls: Model class to validate against

    Returns:
        Validated model instance

    Raises:
        SystemExit: If file loading or validation fails
    """
    try:
        yaml_loader = YAML(typ="safe")
        with open(path, "r") as f:
            data = yaml_loader.load(f)
        obj = model_cls.model_validate(data)
        return obj
    except Exception as e:
        logger.error("Failed to load file", file_path=str(path), error=str(e))
        sys.exit(1)


if __name__ == "__main__":
    sys.exit(app())
