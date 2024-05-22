from moviepy.editor import VideoFileClip
import subprocess


class VideoQuality:
    def __init__(self, label, resolution, video_bitrate, audio_bitrate):
        self.label = label
        self.resolution = resolution
        self.video_bitrate = video_bitrate
        self.audio_bitrate = audio_bitrate

video_qualities = [
    VideoQuality("720p", "1280x720", "5000k", "192k"),
    # VideoQuality("480p", "854x480", "1500k", "128k"),
    # VideoQuality("360p", "640x360", "800k", "96k"),
    # VideoQuality("240p", "426x240", "400k", "64k"),
]


def split_video(input_file, chunk_duration):
    clip = VideoFileClip(input_file)
    total_duration = clip.duration
    start = 0
    chunk_num = 1

    while start < total_duration:
        end = min(start + chunk_duration, total_duration)
        chunk = clip.subclip(start, end)
        chunk_file = f"chunk_{chunk_num}.mp4"
        chunk.write_videofile(chunk_file)
        start = end
        chunk_num += 1

    clip.close()


def encode_chunk(input_file, output_file, resolution, video_bitrate, audio_bitrate):
    cmd = [
        'ffmpeg',
        '-i', input_file,
        '-c:v', 'libx264',
        '-b:v', video_bitrate,
        '-s', resolution,
        '-profile:v', 'main',
        '-level', '3.1',
        '-preset', 'medium',
        '-x264-params', 'scenecut=0:open_gop=0:min-keyint=72:keyint=72',
        '-c:a', 'aac',
        '-b:a', audio_bitrate,
        '-f', 'mp4',
        output_file
    ]
    result = subprocess.run(cmd)
    return result.returncode == 0


def segment_chunk(input_file, output_file):
    cmd = [
        'mp4fragment',
        input_file,
        output_file
    ]
    result = subprocess.run(cmd)
    return result.returncode == 0


def to_mpd(input_file, output_file):
    cmd = [
        'mp4dash',
        input_file,
        '-o', output_file
    ]
    result = subprocess.run(cmd)
    return result.returncode == 0


def is_video(file):
    # TODO: add more video file extensions that are compatible with MPEG_DASH
    allowed_extensions = ['mp4']
    file_extension = file.name.split('.')[-1].lower()
    return file_extension in allowed_extensions

def get_file_path(video_id, quality):
    return f'{video_id}_{quality}.mpd'
