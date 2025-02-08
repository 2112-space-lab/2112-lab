import os
import logging
from flask import Flask
from app.dependencies import Dependencies
from app.routes.health_router import HealthRouter
from app.routes.satellite_router import SatelliteRouter
from app.config import config_class
from app.graphql.schema_loader import load_graphql_schema

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger("propagator-service")

def create_app():
    """Creates and configures the Flask app."""
    flask_app = Flask(__name__)

    deps = Dependencies()
    
    health_router = HealthRouter(deps)
    satellite_router = SatelliteRouter(deps)
    
    flask_app.register_blueprint(health_router.router, url_prefix="/health")
    flask_app.register_blueprint(satellite_router.router, url_prefix="/satellite")

    schema = load_graphql_schema()
    
    @flask_app.route("/graphql", methods=["GET", "POST"])
    def graphql_server():
        from flask import request, jsonify
        from ariadne import graphql_sync
        from ariadne.explorer import ExplorerGraphiQL

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

    return flask_app, schema
