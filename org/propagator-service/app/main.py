import os
import logging
import json
from fastapi import FastAPI
from flask import Flask, Response, jsonify
from ariadne.asgi import GraphQL
from ariadne import make_executable_schema, gql
from dependencies import Dependencies
from routes.health_router import HealthRouter 
from config import config_class
from services.health_resolver import HealthResolver  

# Initialize logger
logger = logging.getLogger("propagator-service")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

# Initialize dependencies
deps = Dependencies()

# Flask REST setup
flask_app = Flask(__name__)

# Register Flask routes
health_router = HealthRouter(deps)  # Initialize health router
flask_app.register_blueprint(health_router.router, url_prefix="/health")

# Read schema directory from environment variable or use default
schema_directory = getattr(config_class, "SCHEMA_DIRECTORY", os.getenv("SCHEMA_DIRECTORY", "./graphql_schemas/"))

# Ensure schema directory exists
if not os.path.exists(schema_directory):
    logger.warning(f"Schema directory not found: {schema_directory}. GraphQL API may not work correctly.")
    type_defs = ""  # Empty schema
else:
    # Load and combine multiple GraphQL schemas
    schema_files = [os.path.join(schema_directory, f) for f in os.listdir(schema_directory) if f.endswith(".graphqls")]

    if not schema_files:
        logger.warning(f"No GraphQL schema files found in: {schema_directory}. GraphQL API may not work correctly.")
        type_defs = ""
    else:
        type_defs = "\n".join(open(schema_file, "r").read() for schema_file in schema_files)

# Initialize GraphQL resolvers
health_resolver = HealthResolver(deps)

# Create GraphQL schema only if we have definitions
if type_defs:
    schema = make_executable_schema(gql(type_defs))
    logger.info("GraphQL schema successfully loaded.")
else:
    schema = None
    logger.warning("GraphQL schema is empty. The GraphQL API may not function as expected.")

# FastAPI GraphQL setup
fastapi_app = FastAPI()

if schema:
    fastapi_app.add_route("/graphql", GraphQL(schema, debug=True))
    logger.info("GraphQL route added at /graphql")
else:
    logger.warning("GraphQL route was not added because the schema is missing.")

# Running the service with Flask
if __name__ == "__main__":
    logger.info("Starting Propagator Service")
    flask_app.run(debug=True, port=5000)
