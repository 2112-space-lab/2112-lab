import base64
import uuid
import json
from dataclasses import asdict, is_dataclass
from datetime import datetime
from typing import Any

from .event_root import EventRoot


class EventBuilder:
    """
    A generic builder class to create an EventRoot from any model.
    """

    def __init__(self, event_type: str, model: Any, comment: str = ""):
        """
        Initialize the EventBuilder with an event type and a model.
        :param event_type: The type of the event (e.g., "PROPAGATION_RESULT").
        :param model: The model to be wrapped in the event.
        :param comment: Optional comments for metadata.
        """
        self.event_type = event_type
        self.model = model
        self.comment = comment

    def build(self) -> EventRoot:
        """
        Creates an EventRoot by serializing the model into a base64-encoded payload.
        """

        if not is_dataclass(self.model):
            raise TypeError(f"Expected a dataclass model, got {type(self.model).__name__}")

        model_dict = asdict(self.model)

        json_data = json.dumps(model_dict).encode("utf-8")

        encoded_payload = base64.b64encode(json_data).decode("utf-8")

        return EventRoot(
            event_time_utc=datetime.utcnow().isoformat() + "Z",
            event_uid=str(uuid.uuid4()),
            event_type=self.event_type,
            comment=self.comment,
            payload=encoded_payload,
        )