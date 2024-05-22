from django.contrib import admin
from django.urls import path
from . import views
from django.conf import settings
from django.conf.urls.static import static

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/upload/', views.upload_video, name="upload_video"),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
