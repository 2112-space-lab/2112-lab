from dataclasses import dataclass, asdict
from typing import Dict, Any
import json
import uuid
from datetime import datetime


@dataclass
class EventRoot:
    """
    Standardized event structure for all propagated messages.
    """
    event_time_utc: str
    event_uid: str
    event_type: str
    comment: str = "" 
    payload: Dict[str, Any] = None

    def to_dict(self) -> Dict[str, Any]:
        return asdict(self)

    @staticmethod
    def generate_event(event_type: str, payload: Dict[str, Any], comment: str = "") -> "EventRoot":
        """
        Generates a new EventRoot instance with a unique ID and timestamp.
        """
        return EventRoot(
            event_time_utc=datetime.utcnow().isoformat() + "Z",
            event_uid=str(uuid.uuid4()),
            event_type=event_type,
            comment=comment,
            payload=payload
        )


def parse_event_root(data: Dict) -> EventRoot:
    return EventRoot(**data)


def serialize_event_root(event_root: EventRoot) -> Dict:
    return event_root.to_dict()
