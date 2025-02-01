from dataclasses import dataclass
from typing import Dict

@dataclass
class SatellitePosition:
    id: str
    name: str
    latitude: float
    longitude: float
    altitude: float
    timestamp: str
    uid: str

def parse_satellite_position(data: Dict) -> SatellitePosition:
    return SatellitePosition(**data)

def serialize_satellite_position(satellite_position: SatellitePosition) -> Dict:
    return {
        "id": satellite_position.id,
        "name": satellite_position.name,
        "latitude": satellite_position.latitude,
        "longitude": satellite_position.longitude,
        "altitude": satellite_position.altitude,
        "timestamp": satellite_position.timestamp,
        "uid": satellite_position.uid
    }
