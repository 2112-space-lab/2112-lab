import os

class Config:
    """Base configuration."""
    SECRET_KEY = os.getenv('SECRET_KEY', 'default-secret-key')
    FLASK_ENV = os.getenv('FLASK_ENV', 'development')
    DEBUG = False
    TESTING = False
    SERVER_NAME = os.getenv('SERVER_NAME', 'localhost:5000')


class DevelopmentConfig(Config):
    """Development configuration."""
    DEBUG = True
    FLASK_ENV = 'development'


class ProductionConfig(Config):
    """Production configuration."""
    FLASK_ENV = 'production'
    SERVER_NAME = os.getenv('SERVER_NAME', 'myapp.com')


class TestingConfig(Config):
    """Testing configuration."""
    TESTING = True
    FLASK_ENV = 'testing'
    SERVER_NAME = None
