import pika
import json
import logging
from typing import Dict, Any
from app.dependencies import Dependencies

logger = logging.getLogger(__name__)

class MessageBroker:
    """
    A message broker for publishing messages to RabbitMQ.
    """

    def __init__(self, dependencies: Dependencies):
        self.dependencies = dependencies
        self.rabbitmq_connection = self._get_rabbitmq_connection()

    def _get_rabbitmq_connection(self):
        """
        Ensures that a RabbitMQ connection is established.
        """
        try:
            connection = self.dependencies.rabbitmq_connection
            if not connection or connection.is_closed:
                logger.info("🔄 Reconnecting to RabbitMQ...")
                connection = self.dependencies.create_rabbitmq_connection()
            return connection
        except Exception as e:
            logger.error(f"❌ Failed to connect to RabbitMQ: {e}")
            return None

    def publish_message(self, routing_key: str, message: Dict[str, Any]):
        """
        Publishes a message to the RabbitMQ queue.

        :param routing_key: The queue name or routing key.
        :param message: The message payload as a dictionary.
        """
        if not self.rabbitmq_connection:
            logger.error("❌ Cannot publish message - No RabbitMQ connection available.")
            return

        try:
            channel = self.rabbitmq_connection.channel()
            channel.queue_declare(queue=routing_key, durable=True)

            message_body = json.dumps(message)
            channel.basic_publish(
                exchange='',
                routing_key=routing_key,
                body=message_body,
                properties=pika.BasicProperties(
                    delivery_mode=2,
                )
            )
            logger.info(f"✅ Successfully published message to {routing_key}")

        except pika.exceptions.AMQPError as e:
            logger.error(f"❌ RabbitMQ error publishing message to {routing_key}: {e}")

        except Exception as e:
            logger.error(f"❌ Unexpected error publishing message to {routing_key}: {e}")

    def close_connection(self):
        """
        Closes the RabbitMQ connection.
        """
        try:
            if self.rabbitmq_connection and self.rabbitmq_connection.is_open:
                self.rabbitmq_connection.close()
                logger.info("🔌 RabbitMQ connection closed.")
        except Exception as e:
            logger.error(f"❌ Error closing RabbitMQ connection: {e}")
