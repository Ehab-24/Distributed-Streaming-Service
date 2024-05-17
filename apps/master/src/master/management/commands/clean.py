from django.core.management.base import BaseCommand
import os
import shutil
from django.conf import settings


class Command(BaseCommand):
    help = "Cleans the media/ folder"

    def handle(self, *args, **options):
        media_dir = settings.BASE_DIR / 'media'
        processed_dir = media_dir / 'processed'
        tmp_dir = media_dir / 'tmp'
        uploads_dir = media_dir / 'uploads'
        shutil.rmtree(media_dir)
        os.mkdir(media_dir)
        os.mkdir(processed_dir)
        os.mkdir(tmp_dir)
        os.mkdir(uploads_dir)
        self.stdout.write(self.style.SUCCESS('Successfully cleaned media/ folder'))
