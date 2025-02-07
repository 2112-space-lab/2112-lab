import json
import base64
import uuid
from datetime import datetime
from typing import Any
from dataclasses import asdict, is_dataclass

from generated.event_root import EventRoot
from generated.enums import EventType
from core.message_broker import MessageBroker

import logging
logger = logging.getLogger(__name__)

class EventEmitter:
    """
    A utility class for building and emitting events.
    """

    def __init__(self, message_broker: MessageBroker):
        self.message_broker = message_broker

    def emit(self, event_type: EventType, model: Any, comment: str = ""):
        """
        Creates an EventRoot and publishes it via RabbitMQ.
        """
        try:
            if not isinstance(event_type, EventType):
                raise ValueError(f"❌ Invalid event type: {event_type}. Must be an instance of EventType Enum.")

            if not is_dataclass(model):
                raise TypeError(f"❌ Expected a dataclass model, got {type(model).__name__}")

            model_dict = asdict(model)
            json_data = json.dumps(model_dict).encode("utf-8")
            encoded_payload = base64.b64encode(json_data).decode("utf-8")

            event = EventRoot(
                event_time_utc=datetime.utcnow().isoformat() + "Z",
                event_uid=str(uuid.uuid4()),
                event_type=event_type,
                comment=comment,
                payload=encoded_payload
            )

            routing_key = f"events.{event_type.value.lower()}"
            event_dict = asdict(event)
            self.message_broker.publish_message(routing_key, event_dict)

            logger.info(f"✅ Event emitted: {event_type.value} | UID: {event.event_uid}")

        except Exception as e:
            logger.error(f"❌ Failed to emit event: {e}")
