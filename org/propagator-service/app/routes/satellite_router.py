import logging
import json
from flask import Blueprint, request, jsonify, Response
from dependencies import Dependencies
from services.satellite_service import propagate_satellite_position
from generated.models import PropagationRequestInput

class SatelliteRouter:
    """
    A dedicated router for handling satellite propagation requests.
    """

    def __init__(self, dependencies: Dependencies):
        """
        Initializes the SatelliteRouter with dependencies.
        """
        self.router = Blueprint("satellite_router", __name__)
        self.dependencies = dependencies
        self.logger = logging.getLogger("satellite-router")

        self._register_routes()

    def _register_routes(self):
        """
        Registers all satellite-related routes.
        """

        @self.router.route("/satellite/propagate", methods=["POST"])
        def propagate():
            """
            Endpoint to propagate satellite positions based on TLE data.
            """
            try:
                request_data = request.get_json()
                self.logger.info(f"Received payload: {json.dumps(request_data, indent=4)}")

                try:
                    propagation_request = PropagationRequestInput(**request_data)
                except Exception as e:
                    self.logger.error(f"Invalid request data: {str(e)}")
                    return jsonify({"error": f"Invalid request data: {str(e)}"}), 400

                self.logger.info(f"Propagating satellite {propagation_request.norad_id} from {propagation_request.start_time}")

                positions = propagate_satellite_position(
                    propagation_request.norad_id,
                    propagation_request.tle_line_1,
                    propagation_request.tle_line_2,
                    propagation_request.start_time,
                    propagation_request.duration_minutes,
                    propagation_request.interval_seconds,
                )

                return Response(json.dumps({"positions": positions}, indent=4), mimetype="application/json"), 200

            except Exception as e:
                self.logger.error(f"Error propagating satellite positions: {str(e)}")
                return jsonify({"error": f"Error propagating satellite positions: {str(e)}"}), 500
