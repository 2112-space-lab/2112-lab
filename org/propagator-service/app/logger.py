import logging
import os
import sys
import json
import traceback
from app.logging.handlers import RotatingFileHandler
from datetime import datetime

try:
    from opentelemetry.trace import get_current_span
except ImportError:
    get_current_span = None


class LokiJSONFormatter(logging.Formatter):
    """Custom JSON log formatter for Grafana Loki / Tempo."""

    def format(self, record):
        log_record = {
            "timestamp": datetime.utcnow().isoformat() + "Z",
            "level": record.levelname,
            "message": record.getMessage(),
            "module": record.module,
            "function": record.funcName,
            "line": record.lineno,
            "service": "propagator-service",
            "hostname": os.getenv("HOSTNAME", "localhost"),
        }

        if get_current_span:
            span = get_current_span()
            if span and span.get_span_context():
                log_record["trace_id"] = f"{span.get_span_context().trace_id:032x}"
                log_record["span_id"] = f"{span.get_span_context().span_id:016x}"

        if record.exc_info:
            log_record["exception"] = "".join(traceback.format_exception(*record.exc_info))

        return json.dumps(log_record)


def setup_logger():
    """Configures structured JSON logging for Grafana Loki / Tempo."""
    log_level = os.getenv("LOG_LEVEL", "INFO").upper()

    logger = logging.getLogger("propagator-service")
    logger.setLevel(getattr(logging, log_level, logging.INFO))

    if logger.hasHandlers():
        logger.handlers.clear()

    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(LokiJSONFormatter())
    console_handler.setLevel(getattr(logging, log_level, logging.INFO))

    log_dir = "propagator-service/logs"
    os.makedirs(log_dir, exist_ok=True)
    file_handler = RotatingFileHandler(
        f"{log_dir}/log", maxBytes=5 * 1024 * 1024, backupCount=3
    )
    file_handler.setFormatter(LokiJSONFormatter())
    file_handler.setLevel(getattr(logging, log_level, logging.INFO))

    logger.addHandler(console_handler)
    logger.addHandler(file_handler)

    return logger


logger = setup_logger()

logger.info("Logger is successfully configured and writing to both stdout and file.")
