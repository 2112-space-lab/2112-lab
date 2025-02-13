import logging
from app.dependencies import Dependencies
from app.core.event_emitter import EventEmitter
from app.models.generated.event_root import EventRoot

logger = logging.getLogger(__name__)

class EventHandler(Dependencies, EventEmitter):
    """
    Base handler class that processes incoming events.
    Inherit from this class to implement specific event handling logic.
    """

    def __init__(self, dependencies: Dependencies):
        Dependencies.__init__(self) 
        EventEmitter.__init__(self, dependencies) 

    def run(self, event: EventRoot):
        """
        Override this method in subclasses to handle incoming events.
        """
        raise NotImplementedError("Subclasses must implement the run method.")
