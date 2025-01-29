import logging
from datetime import datetime
from typing import List
from dependencies import Dependencies
from generated.health_status import HealthStatus, DependencyStatus, serialize_health_status

class HealthService:
    """Service to check health status of dependencies like Redis & RabbitMQ."""

    def __init__(self, dependencies: Dependencies):
        self.logger = logging.getLogger("health-service")
        self.dependencies = dependencies

    def get_trace_id(self) -> str:
        """Extracts OpenTelemetry trace ID if available."""
        return self.dependencies.get_trace_id()

    def log_error(self, message: str):
        """Helper function to log errors with trace ID."""
        self.logger.error(message, extra={"trace_id": self.get_trace_id()})

    def check_redis(self) -> DependencyStatus:
        """Check Redis availability using the Dependencies class."""
        if not self.dependencies.redis_client:
            self.log_error("Redis is unavailable")
            return DependencyStatus(name="Redis", status="unavailable", message="Cannot connect to Redis")

        try:
            self.dependencies.redis_client.ping()
            self.logger.info("Redis is available", extra={"trace_id": self.get_trace_id()})
            return DependencyStatus(name="Redis", status="connected", message=None)
        except Exception as e:
            self.log_error(f"Redis health check failed: {str(e)}")
            return DependencyStatus(name="Redis", status="unavailable", message=str(e))

    def check_rabbitmq(self) -> DependencyStatus:
        """Check RabbitMQ availability using the Dependencies class."""
        if not self.dependencies.rabbitmq_connection:
            self.log_error("RabbitMQ is unavailable")
            return DependencyStatus(name="RabbitMQ", status="unavailable", message="Cannot connect to RabbitMQ")

        try:
            channel = self.dependencies.rabbitmq_connection.channel()
            channel.queue_declare(queue="health_check_queue", passive=True)
            self.logger.info("RabbitMQ is available", extra={"trace_id": self.get_trace_id()})
            return DependencyStatus(name="RabbitMQ", status="connected", message=None)
        except Exception as e:
            self.log_error(f"RabbitMQ health check failed: {str(e)}")
            return DependencyStatus(name="RabbitMQ", status="unavailable", message=str(e))

    def get_health_status(self) -> HealthStatus:
        """
        Returns the aggregated health status of dependencies as a HealthStatus object.
        """
        dependencies = [self.check_redis(), self.check_rabbitmq()]
        status = "healthy" if all(dep.status == "connected" for dep in dependencies) else "degraded"

        health_status = HealthStatus(
            serviceName="Propagator Service",
            status=status,
            timestamp=datetime.utcnow().isoformat() + "Z",
            dependencies=dependencies,
            traceId=self.get_trace_id(),
        )

        return health_status
