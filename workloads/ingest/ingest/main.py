from pathlib import Path
import sys
from typing import Dict, Any

import typer # type: ignore

from ruamel.yaml import YAML # type: ignore
from .log import setup_logging
from .sources import source_factory
from .targets import target_factory

app = typer.Typer()

# Module-level constants for default values
CONFIG_PATH = typer.Option(..., '--config', '-c', help='Path to config file')
DEBUG_FLAG = typer.Option(False, '--debug', '-d', help='Enable debug mode')
ENV_FLAG = typer.Option("dev", '--env', '-e', help='Environment to use')

log = setup_logging(debug=DEBUG_FLAG)

@app.command() 
def ingest(
    configPath: Path = CONFIG_PATH,
    debug: bool = DEBUG_FLAG,
    env: str = ENV_FLAG,
):
    """Ingest data from a source database to a target system."""
    log.info("Starting ingestion process")
    setup_logging(debug=debug)

    config = load_config(configPath)
    try:
        # Initialize source and target instances based on the environment
        source = source_factory(config, env)
        target = target_factory(config, env)
        
        # Validate connection
        if not source.validate_connection():
            log.error("Failed to validate source connection")
            sys.exit(1)
        
        log.info(f"Successfully created {source.config.type} source: {source.config.name}")
        if not target.validate_connection():
            log.error("Failed to validate target connection")
            sys.exit(1)
        log.info(f"Successfully created {target.config.type} target: {target.config.name}")
        
        # TODO: Implement actual extraction logic
    except ValueError as e:
        log.error(f"Configuration error: {e}")
        sys.exit(1)
    except Exception as e:
        log.error(f"Unexpected error: {e}")
        sys.exit(1)


def load_config(configPath: Path) -> Dict[str, Any]:
    try:
        log.debug("Loading configuration file: %s", str(configPath))
        yaml_loader = YAML(typ='safe')
        with open(configPath, "r") as f:
            config = yaml_loader.load(f)
        log.debug("Configuration file loaded successfully")
        return config
    except Exception as e:
        log.error("Failed to load configuration file", 
                 config_path=str(configPath), 
                 error=str(e))
        sys.exit(1)

if __name__ == "__main__":
    app()
