from django.db import models
from copy import deepcopy
from django.conf import settings
import random

class ChunkServer(models.Model):
    id = models.AutoField(primary_key=True)
    ip = models.CharField(max_length=100)
    port = models.IntegerField()
    is_active = models.BooleanField(default=True)
    last_heartbeat = models.DateTimeField(auto_now=True)

    @staticmethod
    def get_active():
        servers = ChunkServer.objects.filter(is_active=True)
        return list(servers)

    @staticmethod
    def n_active():
        return ChunkServer.objects.filter(is_active=True).count()

class Video(models.Model):
    id = models.AutoField(primary_key=True)
    title = models.CharField(max_length=100)
    description = models.TextField()
    replication_factor = models.IntegerField()
    chunk_duration = models.IntegerField()
    duration = models.FloatField()


class Chunk(models.Model):
    id = models.AutoField(primary_key=True)
    video = models.ForeignKey(Video, on_delete=models.CASCADE)
    start_time = models.IntegerField()
    end_time = models.IntegerField()
    replicas = models.ManyToManyField(ChunkServer)
    n_replicas = models.IntegerField()

    def assign_replicas(self, servers):
        """
        Assign replicas to this chunk.

        Args:
            servers (list): List of ChunkServer instances to be assigned as replicas.
        """
        for server in servers:
            self.replicas.add(server)


class ChunkCreator:
    def __init__(self, servers, video, chunk_count):
        self.servers = servers
        self.video = video
        self.chunk_count = chunk_count

    def select_chunk_server(self):
        server = random.choice(self.servers)
        self.servers.remove(server)
        return server

    def create_chunks(self, chunk_duration):
        chunks = []
        total_duration = self.video.duration

        for i in range(self.chunk_count):
            chunk = Chunk.objects.create(
                video=self.video,
                start_time=i * chunk_duration,
                end_time=min((i + 1) * chunk_duration, total_duration),
                n_replicas=self.video.replication_factor
            )
            initial_servers = deepcopy(self.servers)
            chunk_servers = [self.select_chunk_server() for _ in range(self.video.replication_factor)]

            self.servers = initial_servers
            chunk.assign_replicas(chunk_servers)
            chunks.append(chunk)
        return chunks

