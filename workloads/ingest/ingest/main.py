from pathlib import Path
import sys

import typer # type: ignore

from ruamel.yaml import YAML # type: ignore
from .log import setup_logging
from .sources import source_factory
from .core import Config, Catalog
app = typer.Typer()

# Module-level constants for default values
CONFIG_PATH = typer.Option(..., '--config', '-c', help='Path to config file')
CATALOG_PATH = typer.Option(..., '--catalog', '-cat', help='Path to catalog file')
DEBUG_FLAG = typer.Option(False, '--debug', '-d', help='Enable debug mode')
ENV_FLAG = typer.Option("dev", '--env', '-e', help='Environment to use')

logger = setup_logging(debug=DEBUG_FLAG)

@app.command() 
def ingest(
    configPath: Path = CONFIG_PATH,
    catalogPath: Path = CATALOG_PATH,
    debug: bool = DEBUG_FLAG,
    env: str = ENV_FLAG,
):
    """Ingest data from a source database to a target system.
    Args:
        configPath: Path to config file
        debug: Enable debug mode
        env: Environment to use
    """
    # Initialize logger with the debug parameter from the function
    logger = setup_logging(debug=debug)
    logger.info("Starting ingestion process") 

    if env.lower() not in ["dev", "prod"]:
        raise ValueError(f"Invalid Env flag: {env} expected `dev` or `prod`")
    
    config = load_yaml_model(configPath, Config)
    catalog = load_yaml_model(catalogPath, Catalog) 
    
    try:
        # Get all unique sources from the catalog
        sources = catalog.get_sources()
        logger.info(f"Found {len(sources)} sources in catalog: {', '.join(sources)}")
        
        # Process each source and its tables
        for source_name in sources:
            logger.info(f"Processing source: {source_name}")
            
            # Get all tables for this source
            tables = catalog.get_tables_by_source(source_name)
            logger.info(f"Found {len(tables)} tables for source: {source_name}")
            
            # Initialize source based on the catalog's source configuration
            source = source_factory(config, source_name, env)
            
            # Validate connection
            logger.debug(f"Validating connection to source: {source.config.name}")
            if not source.validate_connection():
                raise ValueError(f"Failed to validate source connection: {source.config.name}")
            logger.info(f"Successfully connected to database: `{source.config.name}`")

            # Extract data using streaming for all tables from this source
            for table in tables:
                logger.info(f"Starting extraction of table: {table.name} from source: {source_name}")
                chunk_size = config.params.chunk_size
                row_count = 0
                for row in source.extract(table=table, chunk_size=chunk_size):
                    # TODO: Add logging for row processing
                    row_count += 1
                    if row_count % 1000 == 0:
                        logger.info(f"Processed {row_count} rows from table: {table.name}")
                
                logger.info(f"Completed extraction of {row_count} rows from table: {table.name}")
            
            # Close the source connection after processing all its tables
            source.close()
            logger.info(f"Completed processing source: {source_name}")
        
    except ValueError as e:
        logger.error(f"Configuration error: {e}")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        sys.exit(1)


def load_yaml_model(path: Path, model_cls):
    try:
        logger.debug("Loading file: %s", str(path))
        yaml_loader = YAML(typ='safe')
        with open(path, "r") as f:
            data = yaml_loader.load(f)
        obj = model_cls.model_validate(data)
        logger.info("File loaded and validated successfully: %s", str(path))
        return obj
    except Exception as e:
        logger.error("Failed to load file", file_path=str(path), error=str(e))
        sys.exit(1)


if __name__ == "__main__":
    app()
