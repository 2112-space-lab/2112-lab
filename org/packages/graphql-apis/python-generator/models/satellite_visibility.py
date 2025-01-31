from dataclasses import dataclass
from typing import Dict
from .user_location import UserLocation

@dataclass
class SatelliteVisibility:
    satelliteId: str
    satelliteName: str
    aos: str
    los: str
    userLocation: UserLocation
    uid: str

def parse_satellite_visibility(data: Dict) -> SatelliteVisibility:
    user_location = parse_user_location(data["userLocation"])
    return SatelliteVisibility(
        satelliteId=data["satelliteId"],
        satelliteName=data["satelliteName"],
        aos=data["aos"],
        los=data["los"],
        userLocation=user_location,
        uid=data["uid"]
    )

def serialize_satellite_visibility(satellite_visibility: SatelliteVisibility) -> Dict:
    return {
        "satelliteId": satellite_visibility.satelliteId,
        "satelliteName": satellite_visibility.satelliteName,
        "aos": satellite_visibility.aos,
        "los": satellite_visibility.los,
        "userLocation": serialize_user_location(satellite_visibility.userLocation),
        "uid": satellite_visibility.uid
    }
