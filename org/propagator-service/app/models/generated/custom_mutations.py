# Generated by ariadne-codegen

from typing import Any, Dict, Optional

from .custom_fields import SatellitePositionFields
from .custom_typing_fields import GraphQLField
from .input_types import PropagationRequestInput, UserLocationInput


class Mutation:
    @classmethod
    def request_satellite_visibilities(
        cls, uid: str, user_location: UserLocationInput, start_time: str, end_time: str
    ) -> GraphQLField:
        arguments: Dict[str, Dict[str, Any]] = {
            "uid": {"type": "String!", "value": uid},
            "userLocation": {"type": "UserLocationInput!", "value": user_location},
            "startTime": {"type": "String!", "value": start_time},
            "endTime": {"type": "String!", "value": end_time},
        }
        cleared_arguments = {
            key: value for key, value in arguments.items() if value["value"] is not None
        }
        return GraphQLField(
            field_name="requestSatelliteVisibilities", arguments=cleared_arguments
        )

    @classmethod
    def propagate_satellite_position(
        cls, request: PropagationRequestInput
    ) -> SatellitePositionFields:
        arguments: Dict[str, Dict[str, Any]] = {
            "request": {"type": "PropagationRequestInput!", "value": request}
        }
        cleared_arguments = {
            key: value for key, value in arguments.items() if value["value"] is not None
        }
        return SatellitePositionFields(
            field_name="propagateSatellitePosition", arguments=cleared_arguments
        )
