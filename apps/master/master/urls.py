from django.contrib import admin
from django.urls import path
from . import views

urlpatterns = [
    path('admin/', admin.site.urls),
    path('create/', views.create_video, name='create_video'),
    path('all/', views.get_videos, name='get_videos'),
    path('replica/servers/', views.get_replication_servers, name='replica_servers'),
    path('replica/success/', views.notify_replication, name='notify_replication'),
]
