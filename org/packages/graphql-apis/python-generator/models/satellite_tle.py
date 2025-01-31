from dataclasses import dataclass
from typing import Dict

@dataclass
class SatelliteTle:
    id: str
    name: str
    tleLine1: str
    tleLine2: str
    uid: str

def parse_satellite_tle(data: Dict) -> SatelliteTle:
    return SatelliteTle(**data)

def serialize_satellite_tle(satellite_tle: SatelliteTle) -> Dict:
    return {
        "id": satellite_tle.id,
        "name": satellite_tle.name,
        "tleLine1": satellite_tle.tleLine1,
        "tleLine2": satellite_tle.tleLine2,
        "uid": satellite_tle.uid
    }
