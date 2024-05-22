# Generated by Django 4.2.13 on 2024-05-17 14:57

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('master', '0001_initial'),
    ]

    operations = [
        migrations.CreateModel(
            name='VideoMetaData',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('resolution', models.CharField(max_length=10)),
                ('video_bitrate', models.CharField(max_length=10)),
                ('audio_bitrate', models.CharField(max_length=10)),
                ('file_path', models.CharField(max_length=1024)),
                ('uploaded_at', models.DateTimeField(auto_now_add=True)),
                ('title', models.CharField(max_length=100)),
            ],
        ),
        migrations.RemoveField(
            model_name='video',
            name='title',
        ),
        migrations.RemoveField(
            model_name='video',
            name='uploaded_at',
        ),
        migrations.AlterField(
            model_name='video',
            name='file',
            field=models.FileField(upload_to='videos/'),
        ),
    ]