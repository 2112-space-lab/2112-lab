import json
import logging
import dataclasses
import uuid
from skyfield.api import EarthSatellite, load, wgs84
from datetime import datetime, timedelta
from typing import List

from dependencies import Dependencies
from generated.propagation_request import PropagationRequestInput
from generated.satellite_position import SatellitePosition
from generated.satellite_tle_propagated import SatelliteTlePropagated
from generated.enums import EventType
from core.message_broker import MessageBroker
from core.event_builder import EventBuilder
from core.event_emitter import EventEmitter

logger = logging.getLogger(__name__)

class Propagator:
    def __init__(self, dependencies: Dependencies):
        self.dependencies = dependencies
        self.message_broker = MessageBroker(dependencies)
        self.redis_client = dependencies.redis_client
        self.event_emitter = EventEmitter(self.message_broker) 

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

    def propagate(self, request: PropagationRequestInput) -> str:
        """
        Propagate satellite positions based on TLE data and store in Redis.
        """
        try:
            start_time = self.normalize_and_parse_iso_date(request.start_time)

            ts = load.timescale()
            satellite = EarthSatellite(request.tle_line_1, request.tle_line_2, str(request.norad_id), ts)

            end_time = start_time + timedelta(minutes=request.duration_minutes)
            current_time = start_time
            positions = []

            store_key = f"satellite:{request.norad_id}:positions:{uuid.uuid4()}"

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
                    id=str(request.norad_id),
                    name=f"Satellite {request.norad_id}",
                    latitude=subpoint.latitude.degrees,
                    longitude=subpoint.longitude.degrees,
                    altitude=subpoint.elevation.km,
                    timestamp=current_time.isoformat(),
                )
                
                positions.append(dataclasses.asdict(position))
                current_time += timedelta(seconds=request.interval_seconds)

            self.store_positions_in_redis(store_key, positions)

            self.publish_propagation_event(request, store_key)

            return store_key

        except Exception as e:
            logger.error(f"❌ Error propagating satellite position: {e}")
            raise ValueError(f"Error propagating satellite position: {e}")

    def store_positions_in_redis(self, store_key: str, positions: List[dict]):
        """
        Stores satellite positions in Redis under a unique key.
        """
        try:
            json_data = json.dumps(positions)
            self.redis_client.set(store_key, json_data)
            self.redis_client.expire(store_key, 86400)
            logger.info(f"✅ Stored positions in Redis with key: {store_key}")
        except Exception as e:
            logger.error(f"❌ Failed to store positions in Redis: {e}")

    def publish_propagation_event(self, request: PropagationRequestInput, store_key: str):
        """
        Publishes an event with the stored key instead of raw positions.
        """
        try:
            propagated_event = SatelliteTlePropagated(
                norad_id=request.norad_id,
                tle_line_1=request.tle_line_1,
                tle_line_2=request.tle_line_2,
                time_interval=request.interval_seconds,
                store_key=store_key
            )

            self.event_emitter.emit(event_type=EventType.SATELLITE_TLE_PROPAGATED, model=propagated_event)
            logger.info(f"✅ Published propagation event to RabbitMQ with key: {store_key}")

        except Exception as e:
            logger.error(f"❌ Failed to publish propagation event: {e}")
