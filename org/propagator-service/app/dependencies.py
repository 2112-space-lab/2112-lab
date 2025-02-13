import os
import logging
import redis
import pika

try:
    from opentelemetry.trace import get_current_span
except ImportError:
    get_current_span = None

class Dependencies:

    def __init__(self):
        self.logger = self.setup_logger()
        self.redis_client = self.setup_redis()
        self.rabbitmq_url = self.get_rabbitmq_url()
        self.rabbitmq_connection = None

    def setup_logger(self):
        logger = logging.getLogger("propagator-service")
        log_level = os.getenv("LOG_LEVEL", "INFO").upper()
        logger.setLevel(getattr(logging, log_level, logging.INFO))

        handler = logging.StreamHandler()
        handler.setFormatter(logging.Formatter("%(asctime)s - %(levelname)s - %(message)s"))
        logger.addHandler(handler)

        return logger

    def setup_redis(self):
        redis_url = os.getenv("REDIS_URL", "redis://localhost:6379/0")
        try:
            client = redis.Redis.from_url(redis_url, decode_responses=True)
            client.ping()
            self.logger.info("✅ Connected to Redis")
            return client
        except redis.RedisError as e:
            self.logger.error(f"❌ Failed to connect to Redis: {str(e)}")
            return None

    def get_rabbitmq_url(self):
        """Fetch RabbitMQ URL from environment variables."""
        return os.getenv(
            "RABBITMQ_URL",
            f"amqp://{os.getenv('RABBITMQ_USER', '2112')}:{os.getenv('RABBITMQ_PASSWORD', '2112')}@"
            f"{os.getenv('RABBITMQ_HOST', 'localhost')}:{os.getenv('RABBITMQ_PORT', 5672)}/"
        )

    def get_rabbitmq_connection(self):
        """Creates a single RabbitMQ connection if not already created."""
        if self.rabbitmq_connection and not self.rabbitmq_connection.is_closed:
            return self.rabbitmq_connection  # Return existing connection

        try:
            params = pika.URLParameters(self.rabbitmq_url)
            self.rabbitmq_connection = pika.BlockingConnection(params)
            self.logger.info("✅ Connected to RabbitMQ")
            return self.rabbitmq_connection
        except pika.exceptions.AMQPConnectionError as e:
            self.logger.error(f"❌ Failed to connect to RabbitMQ: {str(e)}")
            return None

    def get_trace_id(self):
        if get_current_span:
            span = get_current_span()
            if span and span.get_span_context():
                return hex(span.get_span_context().trace_id)[2:]
        return None
