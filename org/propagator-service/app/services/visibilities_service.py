import json
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List, Dict
from dependencies import Dependencies
from models import VisibilityEvent
import logging

logger = logging.getLogger(__name__)

def compute_single_visibility(
    satellite_id: str,
    satellite_name: str,
    tle_line1: str,
    tle_line2: str,
    start_time: str,
    end_time: str,
    user_location: Dict,
    user_uid: str,
    interval_seconds: int = 10
) -> Dict:
    try:
        start_time = datetime.fromisoformat(start_time.replace("Z", "+00:00"))
        end_time = datetime.fromisoformat(end_time.replace("Z", "+00:00"))
        ts = load.timescale()
        satellite = EarthSatellite(tle_line1, tle_line2, satellite_id, ts)

        user_lat = user_location["latitude"]
        user_lon = user_location["longitude"]
        user_alt = user_location.get("altitude", 0)
        horizon = user_location.get("horizon", 30)

        user_position = wgs84.latlon(user_lat, user_lon, user_alt)
        current_time = start_time

        aos = None
        los = None
        visible = False

        while current_time <= end_time:
            t = ts.utc(current_time.year, current_time.month, current_time.day,
                       current_time.hour, current_time.minute, current_time.second)
            geocentric = satellite.at(t)
            difference = satellite - user_position
            topocentric = difference.at(t)
            alt, _, _ = topocentric.altaz()

            if alt.degrees > horizon:
                if not visible:
                    aos = current_time
                    visible = True
            elif visible:
                los = current_time
                visible = False
                break

            current_time += timedelta(seconds=interval_seconds)

        if aos and los:
            return {
                "satelliteId": satellite_id,
                "satelliteName": satellite_name,
                "aos": aos.isoformat(),
                "los": los.isoformat(),
                "userLocation": user_location,
                "uid": user_uid,
            }
        return None

    except Exception as e:
        logger.error(f"Error computing single visibility: {e}")
        return None
