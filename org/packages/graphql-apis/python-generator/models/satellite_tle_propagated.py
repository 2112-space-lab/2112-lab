from dataclasses import dataclass, asdict
from typing import Dict
from .event_root import EventRoot


@dataclass
class SatelliteTlePropagated:
    """
    Represents a TLE propagation event with time interval and event metadata.
    """
    satellited_id: str
    tle_line_1: str 
    tle_line_2: str 
    start_time_utc: str
    redis_key: str
    duration_minutes: int
    interval_seconds: int

    def to_dict(self) -> Dict:
        return asdict(self)


def parse_satellite_tle_propagated(data: Dict) -> SatelliteTlePropagated:
    return SatelliteTlePropagated(**data)


def serialize_satellite_tle_propagated(event: SatelliteTlePropagated) -> Dict:
    return event.to_dict()
