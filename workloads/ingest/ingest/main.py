from pathlib import Path
import sys

import typer # type: ignore

from ruamel.yaml import YAML # type: ignore
from .log import setup_logging
from .sources import source_factory
from .core.config import Config

app = typer.Typer()

# Module-level constants for default values
CONFIG_PATH = typer.Option(..., '--config', '-c', help='Path to config file')
DEBUG_FLAG = typer.Option(False, '--debug', '-d', help='Enable debug mode')
ENV_FLAG = typer.Option("dev", '--env', '-e', help='Environment to use')

logger = setup_logging(debug=DEBUG_FLAG)

@app.command() 
def ingest(
    configPath: Path = CONFIG_PATH,
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
    
    config = load_config(configPath)
    try:
        # Initialize source based on the environment
        source = source_factory(config, env)
        
        # Validate connection
        if not source.validate_connection():
            logger.error("Failed to validate source connection")
            sys.exit(1)
        logger.info(f"Successfully connected to database: `{source.config.name}`")
        
    except ValueError as e:
        logger.error(f"Configuration error: {e}")
        sys.exit(1)
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        sys.exit(1)


def load_config(configPath: Path) -> Config:
    try:
        logger.debug("Loading configuration file: %s", str(configPath))
        yaml_loader = YAML(typ='safe')
        with open(configPath, "r") as f:
            config_dict = yaml_loader.load(f)
        # Validate and create IngestConfig instance
        config = Config.model_validate(config_dict)
        logger.info("Configuration file loaded and validated successfully")
        return config
    except Exception as e:
        logger.error("Failed to load configuration file", 
                 config_path=str(configPath), 
                 error=str(e))
        sys.exit(1)

if __name__ == "__main__":
    app()
