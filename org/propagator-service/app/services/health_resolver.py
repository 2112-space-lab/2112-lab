from ariadne import QueryType
from services.health_service import HealthService
from dependencies import Dependencies
from generated.health_status import serialize_health_status

query = QueryType()

class HealthResolver:

    def __init__(self, dependencies: Dependencies):
        self.health_service = HealthService(dependencies)

    @query.field("ping")
    def resolve_ping(self, info) -> str:
        return "pong"

    @query.field("ready")
    def resolve_ready(self, info):
        health_status = self.health_service.get_health_status()
        return serialize_health_status(health_status)
