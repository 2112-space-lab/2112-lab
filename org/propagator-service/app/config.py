import os
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

RABBITMQ_HOST = os.getenv("RABBITMQ_HOST", "rabbitmq-service")
RABBITMQ_PORT = int(os.getenv("RABBITMQ_PORT", 5672))
RABBITMQ_USER = os.getenv("RABBITMQ_USER", "guest")
RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD", "guest")
RABBITMQ_QUEUE = os.getenv("RABBITMQ_QUEUE", "satellite_positions")

class Config:
    SECRET_KEY = os.getenv("SECRET_KEY", "default-secret-key")
    FLASK_ENV = os.getenv("FLASK_ENV", "development")
    DEBUG = False
    TESTING = False
    SERVER_NAME = os.getenv("SERVER_NAME", "localhost:5000")

class DevelopmentConfig(Config):
    DEBUG = True
    FLASK_ENV = "development"
    SCHEMA_DIRECTORY = os.getenv("SCHEMA_DIRECTORY", "./propagator-service/app/graphql_schemas/")

class ProductionConfig(Config):
    FLASK_ENV = "production"
    SERVER_NAME = os.getenv("SERVER_NAME", "myapp.com")

class TestingConfig(Config):
    TESTING = True
    FLASK_ENV = "testing"
    SERVER_NAME = None

config_class = {
    "development": DevelopmentConfig,
    "production": ProductionConfig,
    "testing": TestingConfig,
}.get(os.getenv("FLASK_ENV", "development"), Config)
