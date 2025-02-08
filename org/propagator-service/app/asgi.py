from asgiref.wsgi import WsgiToAsgi
from app.factory import create_app

flask_app, schema = create_app()
application = WsgiToAsgi(flask_app)
