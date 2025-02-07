import os
import logging
import json
from flask import Flask, request, jsonify
from asgiref.wsgi import WsgiToAsgi
from ariadne import graphql_sync, make_executable_schema, gql
from ariadne.explorer import ExplorerGraphiQL
from dependencies import Dependencies
from routes.health_router import HealthRouter
from routes.satellite_router import SatelliteRouter
from config import config_class
from services.health_resolver import HealthResolver

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger("propagator-service")

deps = Dependencies()

flask_app = Flask(__name__)
app = flask_app

health_router = HealthRouter(deps)
satellite_router = SatelliteRouter(deps)
flask_register_blueprint(health_router.router, url_prefix="/health")
flask_register_blueprint(satellite_router.router, url_prefix="/satellite")

schema_directory = getattr(config_class, "SCHEMA_DIRECTORY", os.getenv("SCHEMA_DIRECTORY", "./app/graphql_schemas/"))
type_defs = ""

if not os.path.exists(schema_directory):
    logger.warning(f"⚠️ Schema directory not found: {schema_directory}. GraphQL API may not work correctly.")
else:
    schema_files = [os.path.join(schema_directory, f) for f in os.listdir(schema_directory) if f.endswith(".graphqls")]
    if not schema_files:
        logger.warning(f"⚠️ No GraphQL schema files found in: {schema_directory}. GraphQL API may not work correctly.")
    else:
        type_defs = "\n".join(open(schema_file, "r").read() for schema_file in schema_files)

health_resolver = HealthResolver(deps)
if type_defs:
    schema = make_executable_schema(gql(type_defs))
    logger.info("✅ GraphQL schema successfully loaded.")
else:
    schema = None
    logger.warning("⚠️ GraphQL schema is empty. The GraphQL API may not function as expected.")

@flask_route("/graphql", methods=["GET", "POST"])
def graphql_server():
    if request.method == "GET":
        return ExplorerGraphiQL().html(None)

    try:
        data = request.get_json()
        success, result = graphql_sync(schema, data, context_value={"request": request})
        status_code = 200 if success else 400
        return jsonify(result), status_code
    except Exception as e:
        logger.error(f"❌ GraphQL request failed: {str(e)}")
        return jsonify({"error": "GraphQL request failed"}), 500

app = WsgiToAsgi(flask_app)

if __name__ == "__main__":
    import uvicorn
    logger.info("🚀 Starting Propagator Service on http://0.0.0.0:5000")
    uvicorn.run(app, host="0.0.0.0", port=5000)
