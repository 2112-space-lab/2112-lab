import os
import logging
from ariadne import make_executable_schema, gql
from app.config import config_class

logger = logging.getLogger("propagator-service")

def load_graphql_schema():
    """Loads GraphQL schema from files using the configured SCHEMA_DIRECTORY."""
    
    schema_directory = os.path.abspath(config_class.SCHEMA_DIRECTORY)

    logger.info(f"🔍 Current Working Directory (PWD): {os.getcwd()}")
    logger.info(f"📂 Looking for GraphQL schemas in: {schema_directory}")

    type_defs = ""

    if not os.path.exists(schema_directory):
        logger.warning(f"⚠️ Schema directory not found: {schema_directory}. GraphQL API may not work correctly.")
    else:
        schema_files = [os.path.join(schema_directory, f) for f in os.listdir(schema_directory) if f.endswith(".graphqls")]
        if not schema_files:
            logger.warning(f"⚠️ No GraphQL schema files found in: {schema_directory}. GraphQL API may not work correctly.")
        else:
            logger.info(f"📄 Found {len(schema_files)} schema files: {schema_files}")
            type_defs = "\n".join(open(schema_file, "r").read() for schema_file in schema_files)

    if type_defs.strip():
        schema = make_executable_schema(gql(type_defs))
        logger.info("✅ GraphQL schema successfully loaded.")
        return schema
    else:
        logger.warning("⚠️ GraphQL schema files exist but are empty.")
        return None
