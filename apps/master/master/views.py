from .models import ChunkServer, Video, Chunk, ChunkCreator
import random
from .serializers import VideoSerializer
from rest_framework import status
from rest_framework.decorators import api_view
from rest_framework.response import Response
from django.conf import settings
from django.core import serializers


@api_view(['POST'])
def create_video(request):
    """
    Used by the client to set up metadata for a new video.
    """
    serializer = VideoSerializer(data=request.data)
    if not serializer.is_valid():
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
    if serializer.validated_data['replication_factor'] > ChunkServer.n_active():
        return Response({"error": "Replication factor cannot be greater than the number of active servers."}, status=status.HTTP_400_BAD_REQUEST)
    serializer.save()

    video = Video.objects.get(id=serializer.data['id'])
    total_duration = video.duration
    chunk_count = int(total_duration // settings.CHUNK_DURATION + 1)
    servers = ChunkServer.get_active()
    chunk_creator = ChunkCreator(servers, video, chunk_count)
    created_chunks = chunk_creator.create_chunks()

    chunks = []
    for chunk in created_chunks:
        server = random.choice(list(chunk.replicas.all()))
        chunks.append({
            "chunk_id": chunk.id,
            "server": {
                "id": server.id,
                "ip": server.ip,
                "port": server.port,
            }
        })
    response_data = {
        'video_id': serializer.data['id'],
        'chunks': chunks
    }
    return Response(response_data, status=status.HTTP_201_CREATED)


@api_view(['GET'])
def get_replication_servers(request):
    """
    Used by the primary hcunk server to get the list of replica servers for a chunk.
    """
    video_id = request.query_params.get('video_id')
    chunk_id = request.query_params.get('chunk_id')
    chunk = Chunk.objects.get(id=chunk_id, video_id=video_id)
    servers = chunk.replicas.all()
    return Response({
        "servers": [{"server_id": server.id, "server_ip": server.ip, "server_port": server.port} for server in servers]
    })


@api_view(['POST'])
def notify_replication(request):
    """
    Used by the primary chunk server to notify the master server that a chunk has been successfully replicated.
    """
    n_replicas = request.query_params.get('n_replicas')
    video_id = request.query_params.get('video_id')
    chunk_id = request.query_params.get('chunk_id')
    chunk = Chunk.objects.get(id=chunk_id, video_id=video_id)
    chunk.n_replicas = n_replicas
    chunk.save()
    return Response(status=status.HTTP_200_OK)


@api_view(["Get"])
def get_videos(request):
    """
    Returns a list of all videos, and the chunks associated with each video.
    """
    videos = Video.objects.all()
    response_data = []
    for video in videos:
        chunks = Chunk.objects.filter(video=video)
        response_data.append({
            "video": serializers.serialize('json', [video]),
            "chunks": [{"chunk_id": chunk.id, "n_replicas": chunk.n_replicas} for chunk in chunks]
        })
    return Response(response_data)
