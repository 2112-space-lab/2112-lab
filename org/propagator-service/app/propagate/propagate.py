import json
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List, Dict, Any
from dependencies import Dependencies
from models import SatellitePositionFields
import logging

logger = logging.getLogger(__name__)

class Propagator:
    def __init__(self, dependencies: Dependencies):
        self.dependencies = dependencies

    def normalize_and_parse_iso_date(self, iso_date: str) -> datetime:
        try:
            if iso_date.endswith("Z"):
                iso_date = iso_date[:-1] + "+00:00"

            if "." in iso_date:
                main_part, fractional_and_offset = iso_date.split(".", 1)
                fractional, *offset = fractional_and_offset.split("+", 1)

                if int(fractional[:1]) >= 5:
                    iso_date = f"{main_part}+{'+'.join(offset) if offset else ''}"
                    iso_date = str(datetime.fromisoformat(iso_date) + timedelta(seconds=1))
                else:
                    iso_date = f"{main_part}+{'+'.join(offset) if offset else ''}"

            return datetime.fromisoformat(iso_date)

        except ValueError as e:
            logger.error(f"Error parsing ISO date {iso_date}: {e}")
            raise ValueError(f"Error parsing ISO date {iso_date}: {e}")

    def propagate(self, satellite_id: str, tle_line1: str, tle_line2: str, start_time: str, duration_minutes: int, interval_seconds: int) -> List[Dict[str, Any]]:
        try:
            init_start_time = start_time

            start_time = self.normalize_and_parse_iso_date(start_time)

            ts = load.timescale()
            satellite = EarthSatellite(tle_line1, tle_line2, satellite_id, ts)

            end_time = start_time + timedelta(minutes=duration_minutes)
            current_time = start_time
            positions = []

            while current_time <= end_time:
                if not isinstance(current_time, datetime):
                    raise ValueError(f"current_time is not a datetime object: {type(current_time)}")

                t = ts.utc(current_time.year, current_time.month, current_time.day,
                           current_time.hour, current_time.minute, current_time.second)
                geocentric = satellite.at(t)
                subpoint = wgs84.subpoint(geocentric)
                position = SatellitePositionFields(field_name="satellitePosition")
                position.fields(
                    position.id.alias(satellite_id),
                    position.name.alias(satellite_id),
                    position.latitude.alias(subpoint.latitude.degrees),
                    position.longitude.alias(subpoint.longitude.degrees),
                    position.altitude.alias(subpoint.elevation.km),
                    position.timestamp.alias(current_time.isoformat()),
                )
                positions.append(position)
                current_time += timedelta(seconds=interval_seconds)

            return positions

        except Exception as e:
            raise ValueError(f"Error propagating satellite position: {e}")
