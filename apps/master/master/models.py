from django.db import models
from copy import deepcopy
from django.utils import timezone
import requests
import random

class ChunkServer(models.Model):
    id = models.AutoField(primary_key=True)
    scheme = models.CharField(max_length=10)
    ip = models.CharField(max_length=100)
    port = models.IntegerField()
    inactive_count = models.IntegerField(default=0)
    is_active = models.BooleanField(default=True)
    last_heartbeat = models.DateTimeField(auto_now=True)
    stats = models.JSONField(default=dict)

    @staticmethod
    def get_active():
        servers = ChunkServer.objects.filter(is_active=True)
        return list(servers)

    @staticmethod
    def n_active():
        return ChunkServer.objects.filter(is_active=True).count()

    def url(self):
        return f'{self.scheme}://{self.ip}:{self.port}'

    def send_heartbeat(self):
        response = requests.get(self.url() + '/health/')
        self.last_heartbeat = timezone.now()
        self.save()
        return response

    def normalize(self, value, min_value, max_value):
        if max_value == min_value:
            return 1.0
        return (value - min_value) / (max_value - min_value)

    def compute_score(self, weights, min_max_values):
        cpu_usage = self.normalize(self.stats['cpu']['usage'][0], min_max_values['cpu']['min'], min_max_values['cpu']['max'])
        disk_usage = self.normalize(self.stats['disk']['used'], min_max_values['disk']['min'], min_max_values['disk']['max'])
        memory_usage = self.normalize(self.stats['memory']['used'], min_max_values['memory']['min'], min_max_values['memory']['max'])
        score = (weights['cpu'] * cpu_usage) + (weights['disk'] * disk_usage) + (weights['memory'] * memory_usage)
        return score


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
    checksum = models.CharField(max_length=60)

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
        self.weights = {
            'cpu': 0.5,
            'disk': 0.3,
            'memory': 0.2
        }

    def get_min_max_values(self):
        min_max_values = {
            'cpu': {'min': float('inf'), 'max': float('-inf')},
            'disk': {'min': float('inf'), 'max': float('-inf')},
            'memory': {'min': float('inf'), 'max': float('-inf')}
        }
        
        for server in self.servers:
            cpu_usage = server.stats['cpu']['usage'][0]
            disk_used = server.stats['disk']['used']
            memory_used = server.stats['memory']['used']
            
            if cpu_usage < min_max_values['cpu']['min']:
                min_max_values['cpu']['min'] = cpu_usage
            if cpu_usage > min_max_values['cpu']['max']:
                min_max_values['cpu']['max'] = cpu_usage
            
            if disk_used < min_max_values['disk']['min']:
                min_max_values['disk']['min'] = disk_used
            if disk_used > min_max_values['disk']['max']:
                min_max_values['disk']['max'] = disk_used
            
            if memory_used < min_max_values['memory']['min']:
                min_max_values['memory']['min'] = memory_used
            if memory_used > min_max_values['memory']['max']:
                min_max_values['memory']['max'] = memory_used
        
        return min_max_values

    def select_chunk_servers(self, count):
        min_max_values = self.get_min_max_values()
        server_scores = [(server, server.compute_score(self.weights, min_max_values)) for server in self.servers]
        return [server for server, _ in sorted(server_scores, key=lambda x: x[1], reverse=True)[:count]]

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
            chunk_servers = self.select_chunk_servers(self.video.replication_factor)

            self.servers = initial_servers
            chunk.assign_replicas(chunk_servers)
            chunks.append(chunk)
        return chunks


