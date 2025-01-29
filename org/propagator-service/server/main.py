from fastapi import FastAPI, HTTPException, BackgroundTasks, Request
import redis
import json
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
import threading
import logging
from concurrent.futures import ThreadPoolExecutor
import os
from typing import List, Dict

# FastAPI app instance
app = FastAPI()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Redis client setup
redis_client = redis.StrictRedis(host="redis-service", port=6379, decode_responses=True)

# Function to propagate satellite positions
def normalize_and_parse_iso_date(iso_date: str) -> datetime:
    """
    Normalize and parse an ISO 8601 date string to a datetime object, rounded to the nearest second.

    Args:
        iso_date (str): The ISO 8601 date string to normalize and parse.

    Returns:
        datetime: The parsed datetime object, rounded to the nearest second.

    Raises:
        ValueError: If the date string is not a valid ISO 8601 format.
    """
    try:
        # Normalize the ISO date string
        if iso_date.endswith("Z"):
            iso_date = iso_date[:-1] + "+00:00"

        # Handle fractional seconds (remove or round them)
        if "." in iso_date:
            main_part, fractional_and_offset = iso_date.split(".", 1)
            fractional, *offset = fractional_and_offset.split("+", 1)
            
            # Round fractional seconds to the nearest second
            if int(fractional[:1]) >= 5:  # Check the first digit of the fractional part
                iso_date = f"{main_part}+{'+'.join(offset) if offset else ''}"
                iso_date = str(datetime.fromisoformat(iso_date) + timedelta(seconds=1))
            else:
                iso_date = f"{main_part}+{'+'.join(offset) if offset else ''}"

        # Parse the normalized ISO date string
        return datetime.fromisoformat(iso_date)

    except ValueError as e:
        logger.error(f"Error parsing ISO date {iso_date}: {e}")
        raise ValueError(f"Error parsing ISO date {iso_date}: {e}")

def propagate_satellite_position(
    satellite_id: str,
    tle_line1: str,
    tle_line2: str,
    start_time: str,
    duration_minutes: int,
    interval_seconds: int
) -> List[Dict]:
    try:
        init_start_time = start_time  # Save the original start time for debugging

        # Normalize and parse the start_time
        start_time = normalize_and_parse_iso_date(start_time)

        # Load the satellite TLE data
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
            positions.append({
                "id": satellite_id,
                "timestamp": current_time.isoformat(),
                "latitude": subpoint.latitude.degrees,
                "longitude": subpoint.longitude.degrees,
                "altitude": subpoint.elevation.km
            })
            current_time += timedelta(seconds=interval_seconds)

        return positions

    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Error propagating satellite position: {e}")

# Function to publish satellite positions to Redis
def publish_satellite_positions(satellite_id: str, positions: List[Dict]):
    try:
        for pos in positions:
            key = f"satellite_positions:{satellite_id}"
            timestamp = datetime.fromisoformat(pos['timestamp']).timestamp()
            redis_client.zadd(key, {json.dumps(pos): timestamp})
            redis_client.publish("satellite_positions", json.dumps(pos))

        update_event = {
            "event": "event_satellite_positions_updated",
            "satellite_id": satellite_id,
            "start_time": positions[0]["timestamp"] if positions else None,
            "end_time": positions[-1]["timestamp"] if positions else None,
            "positions_count": len(positions),
        }

        redis_client.publish("event_satellite_positions_updated", json.dumps(update_event))
        logger.info(f"Published satellite update event for {satellite_id}")

    except Exception as e:
        logger.error(f"Error publishing satellite position to Redis: {e}")

# Function to compute AOS/LOS for a single visibility event
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


# Function to subscribe to TLE updates
def subscribe_to_tle_updates():
    pubsub = redis_client.pubsub()
    pubsub.subscribe("satellite_tle_updates")

    logger.info("Subscribed to Redis channel: satellite_tle_updates")

    for message in pubsub.listen():
        if message["type"] == "message":
            try:
                satellite_data = json.loads(message["data"])
                tle_line1 = satellite_data.get("line_1")
                tle_line2 = satellite_data.get("line_2")
                satellite_id = satellite_data.get("id")
                epoch = satellite_data.get("epoch", datetime.utcnow().isoformat() + "Z")

                logger.info(f"Received TLE update for satellite {satellite_id}")
                if not tle_line1 or not tle_line2 or not satellite_id:
                    logger.warning(f"Incomplete TLE data: {satellite_data}")
                    return
                # Calculate altitude and orbital period
                ts = load.timescale()
                satellite = EarthSatellite(tle_line1, tle_line2, satellite_id, ts)
                geocentric = satellite.at(ts.now())
                subpoint = wgs84.subpoint(geocentric)
                altitude_km = subpoint.elevation.km

                if altitude_km < 2000:  # LEO
                    pass_duration_minutes = 1440  # 24 hours
                    num_points = 1440  # 1 point per minute
                elif altitude_km < 35786:  # MEO
                    pass_duration_minutes = 180  # 3 hours
                    num_points = 50  # Approx. 1 point every 3.6 minutes
                else:  # GEO
                    pass_duration_minutes = 1440  # 24 hours
                    num_points = 10  # 1 point every 144 minutes

                interval_seconds = (pass_duration_minutes * 60) // num_points

                logger.info(f"Starting propagation for satellite {satellite_id}")
                positions = propagate_satellite_position(
                    satellite_id, tle_line1, tle_line2, epoch,
                    pass_duration_minutes, interval_seconds
                )
                publish_satellite_positions(satellite_id, positions)
                logger.info(f"Finished propagation and publishing for satellite {satellite_id}")

            except Exception as e:
                logger.error(f"Error processing TLE message: {e}")

# Function to subscribe to user visibility events and set a list of visibilities
def subscribe_to_user_visibility_events():
    pubsub = redis_client.pubsub()
    pubsub.psubscribe("user_visibilities_event:*")

    logger.info("Subscribed to Redis channel: user_visibility_event:*")

    for message in pubsub.listen():
        if message["type"] == "pmessage":  # Pattern-based message
            try:
                         # Log the raw incoming message
                logger.info(f"Incoming visibility event: {message}")

                visibilities = json.loads(message["data"])

                results = []
                for visibility_data in visibilities:
                    satellite_id = visibility_data.get("satelliteID")
                    satellite_name = visibility_data.get("satelliteName")
                    start_time = visibility_data.get("startTime")
                    end_time = visibility_data.get("endTime")
                    tle_line1 = visibility_data.get("tleLine1")
                    tle_line2 = visibility_data.get("tleLine2")
                    user_location = visibility_data.get("userLocation")
                    user_uid = visibility_data.get("userUID")

                    if not satellite_id or not tle_line1 or not tle_line2 or not user_location or not user_uid:
                        logger.warning(f"Incomplete data for visibility event: {visibility_data}")
                        continue

                    visibility = compute_single_visibility(
                        satellite_id, satellite_name, tle_line1, tle_line2,
                        start_time, end_time, user_location, user_uid
                    )

                    if visibility:
                        results.append(visibility)

                # Store the results as a list in Redis
                redis_key = f"satellite_visibilities:{visibilities[0]['userUID']}"
                redis_client.set(redis_key, json.dumps(results))
                logger.info(f"Stored visibility results for UID {visibilities[0]['userUID']}")

            except Exception as e:
                logger.error(f"Error processing visibility events: {e}")


# Start subscriptions on app startup
@app.on_event("startup")
def start_subscriptions():
    # threading.Thread(target=subscribe_to_tle_updates, daemon=True).start()
    threading.Thread(target=subscribe_to_user_visibility_events, daemon=True).start()

# Root endpoint for checking service status
@app.get("/")
def read_root():
    return {"message": "Satellite Propagation Service is running"}

# Health check endpoint
@app.get("/health")
def health_check():
    try:
        redis_client.ping()
        return {"status": "healthy"}
    except Exception as e:
        raise HTTPException(status_code=500, detail="Redis connection failed")

# Endpoint to propagate satellite position based on TLE
@app.post("/satellite/propagate")
async def propagate_endpoint(request: Request) -> Dict[str, List[Dict]]:
    """
    Propagate satellite positions and return them wrapped in a JSON object.
    """
    try:
        # Await the request body to parse JSON
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

        # Call the propagation function
        positions = propagate_satellite_position(
            satellite_id, tle_line1, tle_line2, start_time, duration_minutes, interval_seconds
        )

        # Publish the positions to Redis
        publish_satellite_positions(satellite_id, positions)
        logger.info(f"Published positions for satellite {satellite_id}")

        # Return positions wrapped in a "positions" key
        return {"positions": positions}

    except Exception as e:
        logger.error(f"Error propagating satellite positions: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error propagating satellite positions: {str(e)}")
