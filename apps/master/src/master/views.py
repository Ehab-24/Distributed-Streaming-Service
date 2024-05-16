from rest_framework import status
from rest_framework.decorators import api_view
from rest_framework.response import Response
from .serializers import VideoSerializer
from .models import Video, VideoMetaData
from .utils import encode_chunk, is_video, video_qualities, segment_chunk, to_mpd
from django.conf import settings


@api_view(['POST'])
def upload_video(request):
    if 'file' not in request.data:
        return Response({"error": "No file was uploaded."}, status=status.HTTP_400_BAD_REQUEST)
    uploaded_file = request.data['file']
    if not is_video(uploaded_file):
        return Response({"error": "Only video files are allowed."}, status=status.HTTP_400_BAD_REQUEST)

    serializer = VideoSerializer(data=request.data)
    if not serializer.is_valid():
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)
    serializer.save()
    video_id = serializer.data['id']
    tmp_dir = f'{settings.MEDIA_ROOT}/tmp/'
    video_dir = f'{settings.MEDIA_ROOT}/videos'
    upload_dir = f'{settings.MEDIA_ROOT}/uploads'
    file_title = serializer.data['title']

    uploaded_file_path = f'{upload_dir}/{uploaded_file.name}'
    with open(uploaded_file_path, 'wb+') as destination:
        for chunk in uploaded_file.chunks():
            destination.write(chunk)

    # TODO: add more video file extensions that are compatible with MPEG_DASH
    ext = '.mp4'
    for q in video_qualities:
        fragement_label = f'{video_id}_{q.label}'
        out_file = f'{tmp_dir}/{fragement_label}{ext}'
        if not encode_chunk(uploaded_file_path, out_file, q.resolution, q.video_bitrate, q.audio_bitrate):
            return Response({"error": "Failed to encode video."}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        in_file = out_file
        out_file = out_file.replace(ext, f'_frag{ext}')
        if not segment_chunk(in_file, out_file):
            return Response({"error": "Failed to segment video."}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

        in_file = out_file
        out_file = f'{video_dir}/{fragement_label}'
        if not to_mpd(in_file, out_file):
            return Response({"error": "Failed to generate MPEG-DASH files."}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

        video_metadata = VideoMetaData.objects.create(
            title=file_title,
            file_path=out_file,
            resolution=q.resolution,
            video_bitrate=q.video_bitrate,
            audio_bitrate=q.audio_bitrate
        )

    # TODO: save chunk metadata to mastet

    response_data = { "id": video_id }
    return Response(response_data, status=status.HTTP_201_CREATED)


@api_view(['GET'])
def serve_mpd(request):

    video_id = request.query_params.get('id')
    if not video_id:
        return Response({"error": "No video id was provided."}, status=status.HTTP_400_BAD_REQUEST)

    video_metadata = VideoMetaData.objects.get(id=video_id)
    file_path = video_metadata.file_path
    with open(file_path, 'rb') as f:
        return Response(f.read(), content_type='application/dash+xml')
