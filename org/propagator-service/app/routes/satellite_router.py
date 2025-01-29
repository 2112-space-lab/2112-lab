from fastapi import APIRouter, HTTPException, Request
from services.satellite_service import propagate_satellite_position
from dependencies import Dependencies
from models import SatellitePosition
import logging

logger = logging.getLogger(__name__)

router = APIRouter()

@router.post("/satellite/propagate")
async def propagate_endpoint(request: Request) -> Dict[str, List[Dict]]:
    """
    Propagate satellite positions and return them wrapped in a JSON object.
    """
    try:
        data = await request.json()
        logger.info(f"Received payload: {data}")

        tle_line1 = data.get("tle_line1")
        tle_line2 = data.get("tle_line2")
        start_time = data.get("start_time")
        duration_minutes = data.get("duration_minutes", 90)
        interval_seconds = data.get("interval_seconds", 15)
        satellite_id = data.get("norad_id")

        if not tle_line1 or not tle_line2 or not start_time:
            raise HTTPException(status_code=400, detail="TLE data and start time are required")

        logger.info(f"Propagating satellite {satellite_id} from {start_time}")

        positions = propagate_satellite_position(
            satellite_id, tle_line1, tle_line2, start_time, duration_minutes, interval_seconds
        )

        return {"positions": positions}

    except Exception as e:
        logger.error(f"Error propagating satellite positions: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error propagating satellite positions: {str(e)}")
