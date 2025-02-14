import json
import os
import logging
import asyncio
from typing import Dict, Type, List
import aio_pika
from app.core.message_broker import MessageBroker
from app.dependencies import Dependencies
from app.models.generated.enums import EventType
from app.models.generated.event_root import EventRoot, parse_event_root
from app.core.event_handler import EventHandler

logger = logging.getLogger(__name__)

class EventMonitor:
    """
    A class responsible for subscribing to and processing events sequentially while dispatching to handlers asynchronously.
    """

    def __init__(self, dependencies: Dependencies):
        self.dependencies = dependencies
        self.message_broker = MessageBroker(dependencies)
        self.event_handlers: Dict[EventType, List[Type[EventHandler]]] = {}
        self.queue_name = os.getenv("RABBITMQ_INPUT_QUEUE", "propagator.events.all.input")
        self.event_queue = asyncio.Queue()
        self.rabbitmq_url = dependencies.rabbitmq_url

    def register_handler(self, event_type: EventType, handler_class: Type[EventHandler]):
        """
        Registers a handler class for a specific event type.

        :param event_type: The event type to handle.
        :param handler_class: A class that processes the event.
        """
        if event_type not in self.event_handlers:
            self.event_handlers[event_type] = []
        self.event_handlers[event_type].append(handler_class)
        logger.info(f"✅ Registered handler for event: {event_type.value}")

    async def _consume_messages(self):
        """
        Consumes messages from the queue sequentially and processes them.
        Automatically retries connection on failure.
        """
        while True:
            try:
                logger.info("🔄 Waiting events...")
                connection = await aio_pika.connect_robust(self.rabbitmq_url)
                async with connection:
                    channel = await connection.channel()
                    queue = await channel.declare_queue(self.queue_name, durable=True)

                    logger.info(f"📥 Listening for messages on queue: {self.queue_name}...")

                    async for message in queue:
                        async with message.process():
                            await self.event_queue.put(message.body)

            except Exception as e:
                logger.error(f"❌ Error consuming messages: {e}")
                logger.info("🔄 Reconnecting to RabbitMQ in 5 seconds...")
                await asyncio.sleep(5)

    async def _process_events(self):
        """
        Processes events sequentially but dispatches to handlers asynchronously.
        """
        while True:
            message_body = await self.event_queue.get()
            try:
                message = json.loads(message_body)
                event = parse_event_root(message)
                logger.info(f"🔹 Received event: {event.event_type} | UID: {event.event_uid}")

                handler_classes = self.event_handlers.get(event.event_type, [])
                if handler_classes:
                    for handler_class in handler_classes:
                        handler_instance = handler_class(self.dependencies)
                        asyncio.create_task(handler_instance.run(event))
                else:
                    logger.warning(f"⚠️ No registered handler for event type: {event.event_type}")

            except Exception as e:
                logger.error(f"❌ Error processing event: {e}")

    async def start_monitoring(self):
        """
        Starts listening for events and processes them sequentially.
        Ensures tasks do not cause event loop conflicts.
        """
        logger.info(f"📡 Subscribing to queue: {self.queue_name}")

        try:
            loop = asyncio.get_running_loop()
            if loop.is_running():
                logger.warning("⚠️ Event loop already running, starting tasks directly")
                loop.create_task(self._consume_messages())
                loop.create_task(self._process_events())
            else:
                await asyncio.gather(
                    self._consume_messages(),
                    self._process_events()
                )
        except RuntimeError:
            logger.warning("⚠️ No running event loop found, creating a new one.")
            new_loop = asyncio.new_event_loop()
            asyncio.set_event_loop(new_loop)
            new_loop.run_until_complete(
                asyncio.gather(self._consume_messages(), self._process_events())
            )
