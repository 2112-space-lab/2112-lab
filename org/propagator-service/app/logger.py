import logging
import os
import sys
import json
import traceback
from logging.handlers import RotatingFileHandler
from datetime import datetime

try:
    from opentelemetry.trace import get_current_span
except ImportError:
    get_current_span = None

class LokiJSONFormatter(logging.Formatter):
    """Custom JSON log formatter for Grafana Loki / Tempo."""

    def format(self, record):
        # Standard log fields
        log_record = {
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "level": record.levelname,
            "message": record.getMessage(),
            "module": record.module,
            "function": record.funcName,
            "line": record.lineno,
            "service": "propagator-service",
            "hostname": os.getenv("HOSTNAME", "localhost")
        }

        # Add OpenTelemetry trace context if available
        if get_current_span:
            span = get_current_span()
            if span and span.get_span_context():
                log_record["trace_id"] = span.get_span_context().trace_id
                log_record["span_id"] = span.get_span_context().span_id
        
        # If there's an exception, add stack trace
        if record.exc_info:
            log_record["exception"] = "".join(traceback.format_exception(*record.exc_info))

        return json.dumps(log_record)

def setup_logger():
    """Configures structured JSON logging for Grafana Loki / Tempo."""
    log_level = os.getenv("LOG_LEVEL", "INFO").upper()

    # Create logger
    logger = logging.getLogger("propagator-service")
    logger.setLevel(getattr(logging, log_level, logging.INFO))

    # Console Handler (JSON logs for Loki)
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(LokiJSONFormatter())

    # File Handler (Rotating logs)
    log_dir = "propagator-service/logs"
    os.makedirs(log_dir, exist_ok=True)
    file_handler = RotatingFileHandler(f"{log_dir}/app.log", maxBytes=5 * 1024 * 1024, backupCount=3)
    file_handler.setFormatter(LokiJSONFormatter())

    # Add handlers
    logger.addHandler(console_handler)
    logger.addHandler(file_handler)

    return logger

# Initialize logger
logger = setup_logger()
