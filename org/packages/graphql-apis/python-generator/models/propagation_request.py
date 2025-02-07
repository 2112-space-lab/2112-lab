from dataclasses import dataclass
from typing import Dict

@dataclass
class PropagationRequestInput:
    noradId: str
    tleLine1: str
    tleLine2: str
    startTime: str
    durationMinutes: int = 90
    intervalSeconds: int = 15

def parse_propagation_request_input(data: Dict) -> PropagationRequestInput:
    return PropagationRequestInput(**data)

def serialize_propagation_request_input(propagation_request: PropagationRequestInput) -> Dict:
    return {
        "noradId": propagation_request.noradId,
        "tleLine1": propagation_request.tleLine1,
        "tleLine2": propagation_request.tleLine2,
        "startTime": propagation_request.startTime,
        "durationMinutes": propagation_request.durationMinutes,
        "intervalSeconds": propagation_request.intervalSeconds
    }
