from django.apps import AppConfig
import threading

class MasterConfig(AppConfig):
    name = 'master'

    def ready(self):
        from . import tasks
        thread = threading.Thread(target=tasks.background_task)
        thread.daemon = True
        thread.start()

