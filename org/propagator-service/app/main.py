import os
import logging
import json
from flask import Flask, request, jsonify
from ariadne import graphql_sync, make_executable_schema, gql
from ariadne.explorer import ExplorerGraphiQL
from app.dependencies import Dependencies
from app.routes.health_router import HealthRouter 
from app.routes.satellite_router import SatelliteRouter
from app.config import config_class
from app.services.health_resolver import HealthResolver  

logger = logging.getLogger("propagator-service")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

deps = Dependencies()

flask_app = Flask(__name__)
app = flask_app

health_router = HealthRouter(deps) 
satellite_router = SatelliteRouter(deps) 
flask_app.register_blueprint(health_router.router, url_prefix="/health")
flask_app.register_blueprint(satellite_router.router, url_prefix="/satellite")

schema_directory = getattr(config_class, "SCHEMA_DIRECTORY", os.getenv("SCHEMA_DIRECTORY", "./app/graphql_schemas/"))

if not os.path.exists(schema_directory):
    logger.warning(f"Schema directory not found: {schema_directory}. GraphQL API may not work correctly.")
    type_defs = ""
else:
    schema_files = [os.path.join(schema_directory, f) for f in os.listdir(schema_directory) if f.endswith(".graphqls")]

    if not schema_files:
        logger.warning(f"No GraphQL schema files found in: {schema_directory}. GraphQL API may not work correctly.")
        type_defs = ""
    else:
        type_defs = "\n".join(open(schema_file, "r").read() for schema_file in schema_files)

health_resolver = HealthResolver(deps)

if type_defs:
    schema = make_executable_schema(gql(type_defs))
    logger.info("GraphQL schema successfully loaded.")
else:
    schema = None
    logger.warning("GraphQL schema is empty. The GraphQL API may not function as expected.")

@flask_app.route("/graphql", methods=["GET", "POST"])
def graphql_server():
    if request.method == "GET":
        return ExplorerGraphiQL().html(None)
    
    data = request.get_json()
    success, result = graphql_sync(schema, data, context_value={"request": request})
    status_code = 200 if success else 400
    return jsonify(result), status_code

if __name__ == "__main__":
    logger.info("Starting Propagator Service")
    flask_app.run(debug=True, port=5000)
