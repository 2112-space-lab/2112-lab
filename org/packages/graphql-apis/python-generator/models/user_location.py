from dataclasses import dataclass
from typing import Dict

@dataclass
class UserLocation:
    uid: str
    latitude: float
    longitude: float
    radius: float
    horizon: float

@dataclass
class UserLocationInput:
    uid: str
    latitude: float
    longitude: float
    radius: float
    horizon: float

def parse_user_location(data: Dict) -> UserLocation:
    return UserLocation(**data)

def serialize_user_location(user_location: UserLocation) -> Dict:
    return {
        "uid": user_location.uid,
        "latitude": user_location.latitude,
        "longitude": user_location.longitude,
        "radius": user_location.radius,
        "horizon": user_location.horizon
    }

def parse_user_location_input(data: Dict) -> UserLocationInput:
    return UserLocationInput(**data)

def serialize_user_location_input(user_location_input: UserLocationInput) -> Dict:
    return {
        "uid": user_location_input.uid,
        "latitude": user_location_input.latitude,
        "longitude": user_location_input.longitude,
        "radius": user_location_input.radius,
        "horizon": user_location_input.horizon
    }
