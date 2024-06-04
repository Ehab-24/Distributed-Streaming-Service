import time
from .models import ChunkServer
from . import settings

def background_task():
    while True:
        servers = ChunkServer.get_active()
        
        for server in servers:
            response = server.send_heartbeat()
            if response.status_code != 200:
                server.inactive_count += 1
                if server.inactive_count >= 3:
                    server.is_active = False
                    print(f'Server {server.ip} is inactive')
            else:
                data = response.json()
                server.stats = data
                server.inactive_count = 0
            server.save()
        time.sleep(settings.HEARTBEAT_INTERVAL)
