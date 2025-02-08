import os
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

RABBITMQ_HOST = os.getenv("RABBITMQ_HOST", "localhost")
RABBITMQ_PORT = int(os.getenv("RABBITMQ_PORT", 5672))
RABBITMQ_USER = os.getenv("RABBITMQ_USER", "2112")
RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD", "2112")
RABBITMQ_QUEUE = os.getenv("RABBITMQ_QUEUE", "satellite_positions")

DEFAULT_SCHEMA_DIRECTORY = os.path.abspath("./app/graphql/schemas/")

class Config:
    """Base configuration class."""
    
    SECRET_KEY = os.getenv("SECRET_KEY", "default-secret-key")
    FLASK_ENV = os.getenv("FLASK_ENV", "development")
    DEBUG = False
    TESTING = False
    SERVER_NAME = os.getenv("SERVER_NAME", "localhost:5000")
    SCHEMA_DIRECTORY = os.getenv("SCHEMA_DIRECTORY", DEFAULT_SCHEMA_DIRECTORY)

    @classmethod
    def log_config(cls):
        """Logs important config values for debugging."""
        logger.info(f"🔧 Using configuration: {cls.__name__}")
        logger.info(f"🌍 FLASK_ENV: {cls.FLASK_ENV}")
        logger.info(f"🔑 SECRET_KEY: {'HIDDEN' if cls.SECRET_KEY != 'default-secret-key' else 'DEFAULT'}")
        logger.info(f"📂 Schema Directory: {cls.SCHEMA_DIRECTORY}")
        logger.info(f"🐇 RabbitMQ: {RABBITMQ_USER}@{RABBITMQ_HOST}:{RABBITMQ_PORT}, Queue: {RABBITMQ_QUEUE}")

class DevelopmentConfig(Config):
    """Development environment configuration."""
    DEBUG = True
    FLASK_ENV = "development"

class ProductionConfig(Config):
    """Production environment configuration."""
    FLASK_ENV = "production"
    SERVER_NAME = os.getenv("SERVER_NAME", "test")

class TestingConfig(Config):
    """Testing environment configuration."""
    TESTING = True
    FLASK_ENV = "testing"
    SERVER_NAME = None

FLASK_ENV = os.getenv("FLASK_ENV", "development")
config_class = {
    "development": DevelopmentConfig,
    "production": ProductionConfig,
    "testing": TestingConfig,
}.get(FLASK_ENV, DevelopmentConfig)()

config_class.log_config()
