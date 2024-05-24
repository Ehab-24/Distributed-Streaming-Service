from django.db import models

class ChunkServer(models.Model):
    id = models.AutoField(primary_key=True)
    ip = models.CharField(max_length=100)
    port = models.IntegerField()
    is_active = models.BooleanField(default=True)
    last_heartbeat = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f"{self.id} - {self.ip}:{self.port}"

    @staticmethod
    def get_active():
        return ChunkServer.objects.filter(is_active=True)

class Video(models.Model):
    id = models.AutoField(primary_key=True)
    title = models.CharField(max_length=100)
    description = models.TextField()
    replication_factor = models.IntegerField()
    duration = models.IntegerField()

    def __str__(self):
        return f"{self.id} - {self.title}"

class Chunk(models.Model):
    id = models.AutoField(primary_key=True)
    video = models.ForeignKey(Video, on_delete=models.CASCADE)
    start_time = models.IntegerField()
    end_time = models.IntegerField()
    replicas = models.ManyToManyField(ChunkServer)
    n_replicas = models.IntegerField()

    def __str__(self):
        return f"{self.id} - {self.video_id}"
