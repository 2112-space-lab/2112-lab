import dataclasses
import logging
import uuid
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List, Tuple

from app.models.generated.propagation_request import PropagationRequestInput
from app.models.generated.satellite_position import SatellitePosition

logger = logging.getLogger(__name__)

class Propagator:
    def normalize_and_parse_iso_date(self, iso_date: str) -> datetime:
        """
        Normalize and parse ISO 8601 date string to datetime object.
        """
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
            logger.error(f"❌ Error parsing ISO date {iso_date}: {e}")
            raise ValueError(f"Error parsing ISO date {iso_date}: {e}")

    def propagate(self, request: PropagationRequestInput) -> Tuple[str, List[dict]]:
        """
        Propagates satellite positions based on TLE data.
        Returns: (store_key, list of positions)
        """
        try:
            start_time = self.normalize_and_parse_iso_date(request.start_time)
            ts = load.timescale()
            satellite = EarthSatellite(request.tle_line_1, request.tle_line_2, str(request.space_id), ts)

            end_time = start_time + timedelta(minutes=request.duration_minutes)
            current_time = start_time
            positions = []

            store_key = f"satellite:{request.space_id}:positions:{uuid.uuid4()}"

            while current_time <= end_time:
                if not isinstance(current_time, datetime):
                    raise ValueError(f"❌ current_time is not a datetime object: {type(current_time)}")

                t = ts.utc(
                    current_time.year, current_time.month, current_time.day,
                    current_time.hour, current_time.minute, current_time.second
                )
                geocentric = satellite.at(t)
                subpoint = wgs84.subpoint(geocentric)

                position = SatellitePosition(
                    id=str(request.space_id),
                    name=f"Satellite {request.space_id}",
                    latitude=subpoint.latitude.degrees,
                    longitude=subpoint.longitude.degrees,
                    altitude=subpoint.elevation.km,
                    timestamp=current_time.isoformat(),
                )
                
                positions.append(dataclasses.asdict(position))
                current_time += timedelta(seconds=request.interval_seconds)

            return store_key, positions

        except Exception as e:
            logger.error(f"❌ Error propagating satellite position: {e}")
            raise ValueError(f"Error propagating satellite position: {e}")
