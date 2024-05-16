from django.db import models

class Video(models.Model):
    file = models.FileField(upload_to='videos/')


def VideoMetaData(models.Model):
    resolution = models.CharField(max_length=10)
    video_bitrate = models.CharField(max_length=10)
    audio_bitrate = models.CharField(max_length=10)
    file_path = models.CharField(max_length=1024)
    uploaded_at = models.DateTimeField(auto_now_add=True)
    title = models.CharField(max_length=100)

    def __str__(self):
        return f'{self.title} - {self.resolution} - {self.video_bitrate} - {self.audio_bitrate} - {self.uploaded_at} - {self.file_path}'
