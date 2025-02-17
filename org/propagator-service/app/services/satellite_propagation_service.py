import logging
from app.core.base_service import BaseService
from app.dependencies import Dependencies
from app.models.generated.propagation_request import PropagationRequestInput
from app.models.generated.satellite_tle_propagated import SatelliteTlePropagated
from app.models.generated.enums import EventType
from app.propagate.propagate import Propagator

logger = logging.getLogger(__name__)

class SatellitePropagationService(BaseService):
    def __init__(self, dependencies: Dependencies):
        super().__init__(dependencies)
        self.propagator = Propagator()

    def propagate_and_store(self, request: PropagationRequestInput) -> str:
        """
        Calls Propagator to compute positions, then stores in Redis and publishes an event.
        """
        try:
            redis_key, positions = self.propagator.propagate(request)

            self.store_in_redis(redis_key, positions)

            event = SatelliteTlePropagated(
                satellited_id=request.space_id,
                tle_line_1=request.tle_line_1,
                tle_line_2=request.tle_line_2,
                start_time_utc=request.startTime,
                redis_key=redis_key,
                duration_minutes=request.duration_minutes,
                interval_seconds=request.interval_seconds
            )

            self.publish_event(
                event_type=EventType.SATELLITE_TLE_PROPAGATED,
                model=event
            )

            return redis_key

        except Exception as e:
            logger.error(f"❌ Error in propagate_and_store: {e}", exc_info=True)
            raise
