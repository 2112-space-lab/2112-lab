import os
import logging
import asyncio
import threading
from flask import Flask, request, jsonify
from ariadne import graphql_sync
from ariadne.explorer import ExplorerGraphiQL

from app.dependencies import Dependencies
from app.routes.health_router import HealthRouter
from app.routes.satellite_router import SatelliteRouter
from app.config import config_class
from app.graphql.schema_loader import load_graphql_schema
from app.core.event_monitor import EventMonitor
from app.models.generated.enums import EventType
from app.handlers.tle_propagation_requested import TLEPropagationHandler

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger("propagator-service")

def create_app():
    """Creates and configures the Flask app."""
    flask_app = Flask(__name__)

    deps = Dependencies()
    event_monitor = EventMonitor(deps)
    event_monitor.register_handler(EventType.SATELLITE_TLE_PROPAGATION_REQUESTED, TLEPropagationHandler)

    loop = asyncio.get_event_loop()
    if loop.is_running():
        logger.info("📡 Using existing event loop for event monitor")
        loop.create_task(event_monitor.start_monitoring())
    else:
        logger.info("🆕 Starting new event loop for event monitor")
        asyncio.run(event_monitor.start_monitoring())

    health_router = HealthRouter(deps)
    satellite_router = SatelliteRouter(deps)
    flask_app.register_blueprint(health_router.router, url_prefix="/health")
    flask_app.register_blueprint(satellite_router.router, url_prefix="/satellite")

    schema = load_graphql_schema()

    @flask_app.route("/graphql", methods=["GET", "POST"])
    def graphql_server():
        """GraphQL API endpoint."""
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
