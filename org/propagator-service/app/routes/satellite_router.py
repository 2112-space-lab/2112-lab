import logging
import json
import traceback
from flask import Blueprint, request, jsonify
from app.dependencies import Dependencies
from app.models.generated.input_types import PropagationRequestInput
from app.services.satellite_propagation_service import SatellitePropagationService

class SatelliteRouter:
    """
    Router for handling satellite propagation requests.
    """

    def __init__(self, dependencies: Dependencies):
        self.router = Blueprint("satellite_router", __name__)
        self.dependencies = dependencies
        self.logger = logging.getLogger("satellite-router")
        self.propagation_service = SatellitePropagationService(dependencies)

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
                client_ip = request.remote_addr
                headers = dict(request.headers)
                self.logger.info(
                    f"Received request from {client_ip} | Headers: {json.dumps(headers, indent=4)}"
                )

                try:
                    request_data = request.get_json()
                    if not request_data:
                        self.logger.warning(f"Received empty JSON payload from {client_ip}")
                        return jsonify({"error": "Empty request body"}), 400
                except Exception as e:
                    self.logger.error(f"Malformed JSON from {client_ip}: {str(e)}")
                    return jsonify({"error": "Malformed JSON request"}), 400

                self.logger.info(
                    f"Received payload from {client_ip}: {json.dumps(request_data, indent=4)}"
                )

                try:
                    propagation_request = PropagationRequestInput(**request_data)
                except Exception as e:
                    self.logger.error(f"Invalid request data from {client_ip}: {str(e)}")
                    return jsonify({"error": f"Invalid request data: {str(e)}"}), 400

                self.logger.info(
                    f"Propagating satellite {propagation_request.space_id} from {propagation_request.start_time}"
                )

                store_key = self.propagation_service.propagate_and_store(propagation_request)

                self.logger.info(
                    f"Propagation complete for SPACE ID {propagation_request.space_id} | Stored under key: {store_key}"
                )

                return jsonify({"message": "Propagation successful", "store_key": store_key}), 200

            except Exception as e:
                error_trace = traceback.format_exc()
                self.logger.error(
                    f"Error propagating satellite positions: {str(e)}\n{error_trace}"
                )
                return jsonify({"error": f"Internal server error: {str(e)}"}), 500
