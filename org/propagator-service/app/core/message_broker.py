import pika
import json
import logging
from typing import Dict, Any, Callable
from app.dependencies import Dependencies

logger = logging.getLogger(__name__)

class MessageBroker:
    """
    A message broker for publishing and subscribing to messages in RabbitMQ.
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
                connection = self.dependencies.get_rabbitmq_connection()
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

    def subscribe(self, queue_name: str, callback: Callable[[Dict[str, Any]], None]):
        """
        Subscribes to a RabbitMQ queue and consumes messages.

        :param queue_name: The queue name to consume messages from.
        :param callback: A function that will process the message.
        """
        if not self.rabbitmq_connection:
            logger.error("❌ Cannot subscribe - No RabbitMQ connection available.")
            return

        try:
            channel = self.rabbitmq_connection.channel()
            channel.queue_declare(queue=queue_name, durable=True)

            def on_message(channel, method, properties, body):
                try:
                    message = json.loads(body)
                    callback(message)
                    channel.basic_ack(delivery_tag=method.delivery_tag)
                except Exception as e:
                    logger.error(f"❌ Error processing message from {queue_name}: {e}")

            channel.basic_consume(queue=queue_name, on_message_callback=on_message)
            logger.info(f"📥 Listening for messages on {queue_name}...")
            channel.start_consuming()
        
        except pika.exceptions.AMQPError as e:
            logger.error(f"❌ RabbitMQ error while subscribing to {queue_name}: {e}")

        except Exception as e:
            logger.error(f"❌ Unexpected error while subscribing to {queue_name}: {e}")

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
