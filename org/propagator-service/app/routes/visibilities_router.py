from fastapi import APIRouter, HTTPException, Request
from services.visibility_service import compute_single_visibility
from dependencies import Dependencies
from models import VisibilityEvent
import logging

logger = logging.getLogger(__name__)

router = APIRouter()

@router.post("/visibility/compute")
async def compute_visibility_endpoint(request: Request) -> Dict[str, VisibilityEvent]:
    """
    Compute visibility events for a satellite.
    """
    try:
        data = await request.json()
        logger.info(f"Received payload: {data}")

        satellite_id = data.get("satellite_id")
        satellite_name = data.get("satellite_name")
        tle_line1 = data.get("tle_line1")
        tle_line2 = data.get("tle_line2")
        start_time = data.get("start_time")
        end_time = data.get("end_time")
        user_location = data.get("user_location")
        user_uid = data.get("user_uid")

        if not satellite_id or not tle_line1 or not tle_line2 or not start_time or not end_time or not user_location or not user_uid:
            raise HTTPException(status_code=400, detail="All fields are required")

        logger.info(f"Computing visibility for satellite {satellite_id} from {start_time} to {end_time}")

        visibility = compute_single_visibility(
            satellite_id, satellite_name, tle_line1, tle_line2, start_time, end_time, user_location, user_uid
        )

        if visibility:
            return {"visibility": visibility}
        else:
            raise HTTPException(status_code=404, detail="Visibility event not found")

    except Exception as e:
        logger.error(f"Error computing visibility: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error computing visibility: {str(e)}")
