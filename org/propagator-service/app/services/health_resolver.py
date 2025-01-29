from ariadne import QueryType
from services.health_service import HealthService
from dependencies import Dependencies
from generated.custom_fields import HealthStatusFields

query = QueryType()

class HealthResolver:
    """GraphQL Resolver for Health Checks."""

    def __init__(self, dependencies: Dependencies):
        self.health_service = HealthService(dependencies)

    @query.field("ping")
    def resolve_ping(self, info) -> str:
        """Returns a simple response to verify the service is alive."""
        return "pong"

    @query.field("ready")
    def resolve_ready(self, info) -> HealthStatusFields:
        """Checks the availability of Redis & RabbitMQ."""
        return self.health_service.get_health_status()
