import logging
import json
from flask import Blueprint, request, jsonify, Response
from app.dependencies import Dependencies
from app.propagate.propagate import Propagator
from app.generated.input_types import PropagationRequestInput

class SatelliteRouter:
    """
    Router for handling satellite propagation requests.
    """

    def __init__(self, dependencies: Dependencies):
        self.router = Blueprint("satellite_router", __name__)
        self.dependencies = dependencies
        self.logger = logging.getLogger("satellite-router")
        self.propagator = Propagator(dependencies)

        self._register_routes()

    def _register_routes(self):
        """
        Registers all satellite-related routes.
        """

        @self.router.route("/propagate", methods=["POST"])
        def propagate():
            """
            API endpoint to propagate satellite positions.
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

                positions = self.propagator.propagate(propagation_request)

                return jsonify({"positions": positions}), 200

            except Exception as e:
                self.logger.error(f"Error propagating satellite positions: {str(e)}")
                return jsonify({"error": f"Error propagating satellite positions: {str(e)}"}), 500
