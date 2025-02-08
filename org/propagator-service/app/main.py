import logging
import uvicorn
from app.asgi import application

logger = logging.getLogger("propagator-service")

if __name__ == "__main__":
    logger.info("🚀 Starting Propagator Service on http://0.0.0.0:5000")
    uvicorn.run(application, host="0.0.0.0", port=5000)
