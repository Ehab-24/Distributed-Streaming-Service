from .models import ChunkServer, Video, Chunk
import random
from .serializers import VideoSerializer
from rest_framework import status
from rest_framework.decorators import api_view
from rest_framework.response import Response
from django.conf import settings


@api_view(['GET'])
def create_video(request):
    serializer = VideoSerializer(data=request.dqata)
    if not serializer.is_valid():
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
    serializer.save()

    video = Video.objects.get(id=serializer.data['id'])
    total_duration = video.duration
    chunk_count = total_duration // settings.CHUNK_DURATION + 1
    chunks = []
    servers = ChunkServer.get_active()

    def select_chunk_server():
        server = random.choice(servers)
        servers.remove(server)
        return server

    for i in range(chunk_count):
        chunk = Chunk.objects.create(
            video=video,
            start_time=i * settings.CHUNK_DURATION,
            end_time=min((i + 1) * settings.CHUNK_DURATION, total_duration),
            replicas=[select_chunk_server() for _ in range(video.replication_factor)]
        )
        chunks.append(chunk)

    servers = ChunkServer.get_active()
    response_data = {
        'video_id': serializer.data['id'],
        'chunks': [{"chunk_id": chunk.id, "server": select_chunk_server() } for chunk in chunks],
    }
    return Response(response_data, status=status.HTTP_201_CREATED)


@api_view(['GET'])
def get_replication_servers(request):
    video_id = request.query_params.get('video_id')
    chunk_id = request.query_params.get('chunk_id')
    chunk = Chunk.objects.get(id=chunk_id, video_id=video_id)
    servers = chunk.replicas.all()
    return Response([{"server_id": server.id, "server_ip": server.ip, "server_port": server.port} for server in servers])


@api_view(['POST'])
def notify_replication(request):
    """
    Used by the primary chunk server to notify the master server that a chunk is successfully replicated.
    """
    n_replicas = request.query_params.get('n_replicas')
    video_id = request.query_params.get('video_id')
    chunk_id = request.query_params.get('chunk_id')
    chunk = Chunk.objects.get(id=chunk_id, video_id=video_id)
    chunk.n_replicas = n_replicas
    chunk.save()
    return Response(status=status.HTTP_200_OK)
