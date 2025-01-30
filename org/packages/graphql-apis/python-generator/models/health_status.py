from dataclasses import dataclass
from typing import List, Optional
import json
from enum import Enum

class HealthStatusEnum(str, Enum):
    HEALTHY = "HEALTHY"
    DEGRADED = "DEGRADED"
    UNAVAILABLE = "UNAVAILABLE"
    CONNECTED = "CONNECTED"

@dataclass
class DependencyStatus:
    name: str
    status: HealthStatusEnum
    message: Optional[str] = None

@dataclass
class HealthStatus:
    serviceName: str
    status: HealthStatusEnum
    timestamp: str
    dependencies: List[DependencyStatus]
    traceId: Optional[str] = None

def parse_health_status(data: dict) -> HealthStatus:
    dependencies = [DependencyStatus(**dep) for dep in data["dependencies"]]
    return HealthStatus(
        serviceName=data["serviceName"],
        status=HealthStatusEnum(data["status"]),
        timestamp=data["timestamp"],
        dependencies=dependencies,
        traceId=data.get("traceId"),
    )

def serialize_health_status(health_status: HealthStatus) -> dict:
    return {
        "serviceName": health_status.serviceName,
        "status": health_status.status.value,
        "timestamp": health_status.timestamp,
        "dependencies": [
            {
                "name": dep.name,
                "status": dep.status.value,
                "message": dep.message
            } for dep in health_status.dependencies
        ],
        "traceId": health_status.traceId,
    }
