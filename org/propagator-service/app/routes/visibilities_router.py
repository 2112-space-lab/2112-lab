import logging
import json
from flask import Blueprint, request, jsonify, Response
from app.dependencies import Dependencies
from app.services.visibility_service import compute_single_visibility
from app.models.generated.models import VisibilityEvent

class VisibilityRouter:

    def __init__(self, dependencies: Dependencies):
        self.router = Blueprint("visibility_router", __name__)
        self.dependencies = dependencies
        self.logger = logging.getLogger("visibility-router")

        self._register_routes()

    def _register_routes(self):

        @self.router.route("/visibility/compute", methods=["POST"])
        def compute_visibility():
            try:
                request_data = request.get_json()
                self.logger.info(f"Received payload: {json.dumps(request_data, indent=4)}")

                satellite_id = request_data.get("satellite_id")
                satellite_name = request_data.get("satellite_name")
                tle_line1 = request_data.get("tle_line1")
                tle_line2 = request_data.get("tle_line2")
                start_time = request_data.get("start_time")
                end_time = request_data.get("end_time")
                user_location = request_data.get("user_location")
                user_uid = request_data.get("user_uid")

                if not all([satellite_id, tle_line1, tle_line2, start_time, end_time, user_location, user_uid]):
                    self.logger.error("Missing required fields")
                    return jsonify({"error": "All fields are required"}), 400

                self.logger.info(f"Computing visibility for satellite {satellite_id} from {start_time} to {end_time}")

                visibility = compute_single_visibility(
                    satellite_id, satellite_name, tle_line1, tle_line2, start_time, end_time, user_location, user_uid
                )

                if visibility:
                    return Response(json.dumps({"visibility": visibility.dict()}, indent=4), mimetype="application/json"), 200
                else:
                    self.logger.warning(f"Visibility event not found for satellite {satellite_id}")
                    return jsonify({"error": "Visibility event not found"}), 404

            except Exception as e:
                self.logger.error(f"Error computing visibility: {str(e)}")
                return jsonify({"error": f"Error computing visibility: {str(e)}"}), 500
