import pika
import json
from dependencies import Dependencies
import logging
from typing import List, Dict, Any

logger = logging.getLogger(__name__)

class MessageBroker:
    def __init__(self, dependencies: Dependencies):
        self.dependencies = dependencies
        self.rabbitmq_connection = dependencies.rabbitmq_connection

    def publish_message(self, routing_key: str, message: Dict[str, Any]):
        try:
            channel = self.rabbitmq_connection.channel()
            message_body = json.dumps(message)
            channel.basic_publish(
                exchange='',
                routing_key=routing_key,
                body=message_body,
                properties=pika.BasicProperties(
                    delivery_mode=2,
                )
            )
            logger.debug(f"Published message to {routing_key}")
        except Exception as e:
            logger.error(f"Error publishing message to {routing_key}: {e}")