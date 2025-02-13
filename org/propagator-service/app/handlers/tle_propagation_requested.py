import logging
import json
import asyncio
from app.dependencies import Dependencies
from app.core.event_handler import EventHandler
from app.models.generated.event_root import EventRoot
from app.services.satellite_propagation_service import SatellitePropagationService
from app.models.generated.input_types import PropagationRequestInput

logger = logging.getLogger(__name__)

class TLEPropagationHandler(EventHandler):
    def __init__(self, dependencies: Dependencies):
        super().__init__(dependencies)
        self.propagation_service = SatellitePropagationService(dependencies)

    async def run(self, event: EventRoot):
        """Handles incoming TLE propagation events asynchronously."""
        logger.info(f"🚀 Handling TLE Propagation Event - UID: {event.event_uid}")

        try:
            payload_data = event.payload if isinstance(event.payload, dict) else json.loads(event.payload)
            logger.info(f"📦 Processed Payload: {json.dumps(payload_data, indent=2)}")

            propagation_request = PropagationRequestInput(
                norad_id=payload_data.get("norad_id"),
                tle_line_1=payload_data.get("tle_line1"),
                tle_line_2=payload_data.get("tle_line2"),
                start_time=payload_data.get("start_time"),
                duration_minutes=payload_data.get("duration_minutes", 90),
                interval_seconds=payload_data.get("interval_seconds", 15), 
            )

            logger.info(f"📡 Starting propagation for NORAD ID: {propagation_request.norad_id}")
            store_key = await asyncio.to_thread(self.propagation_service.propagate_and_store, propagation_request)
            logger.info(f"✅ Propagation completed and stored with key: {store_key}")

        except json.JSONDecodeError as e:
            logger.error(f"❌ JSON Parsing Error in TLEPropagationHandler: {e}", exc_info=True)
        except KeyError as e:
            logger.error(f"❌ Missing required field in payload: {e}", exc_info=True)
        except Exception as e:
            logger.error(f"❌ Unexpected Error in TLEPropagationHandler: {e}", exc_info=True)
