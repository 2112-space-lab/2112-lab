import json
import logging
from typing import Any
from app.core.event_emitter import EventEmitter
from app.dependencies import Dependencies

logger = logging.getLogger(__name__)

class BaseService:
    """
    Base service providing shared functionalities: Redis storage & event publishing.
    """

    def __init__(self, dependencies: Dependencies):
        self.redis_client = dependencies.redis_client
        self.event_emitter = EventEmitter(dependencies)

    def store_in_redis(self, store_key: str, data: Any, expiration: int = 86400):
        """
        Stores given data in Redis under a unique key with an expiration time.

        :param store_key: The Redis key to store data under.
        :param data: The data to be stored (dict, list, or serializable object).
        :param expiration: Time in seconds before the key expires (default: 86400s = 24h).
        """
        try:
            json_data = json.dumps(data)
            self.redis_client.set(store_key, json_data)
            self.redis_client.expire(store_key, expiration)
            logger.info(f"✅ Stored data in Redis with key: {store_key}")
        except Exception as e:
            logger.error(f"❌ Failed to store data in Redis: {e}")

    def publish_event(self, event_type: str, model: Any, comment: str = ""):
        """
        Publishes an event using EventEmitter.

        :param event_type: The event type (used as routing key).
        :param model: The event payload as a dictionary.
        :param comment: Optional comment describing the event.
        """
        try:
            self.event_emitter.emit(event_type=event_type, model=model, comment=comment)
            logger.info(f"✅ Event published: {event_type}")
        except Exception as e:
            logger.error(f"❌ Failed to publish event: {e}")
