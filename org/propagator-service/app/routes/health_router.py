import json
import logging
from flask import Blueprint, Response, jsonify
from app.dependencies import Dependencies
from app.services.health_service import HealthService
from app.generated.health_status import serialize_health_status

class HealthRouter:
    def __init__(self, dependencies: Dependencies):
        self.router = Blueprint("health_router", __name__)
        self.health_service = HealthService(dependencies)
        self.logger = logging.getLogger("health-router")
        self.register_routes()

    def register_routes(self):
        @self.router.route("/ready", methods=["GET"])
        def ready():
            try:
                status = self.health_service.get_health_status()
                serialized_status = serialize_health_status(status)
                return Response(json.dumps(serialized_status, indent=4), mimetype="application/json")
            except Exception as e:
                self.logger.error(f"Error in /ready route: {str(e)}")
                return jsonify({"error": "Internal Server Error"}), 500

        @self.router.route("/live", methods=["GET"])
        def live():
            return jsonify({"status": "alive"})
