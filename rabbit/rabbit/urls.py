from django.contrib import admin
from django.conf.urls import url, include
from django.urls import path

from rest_framework.authtoken import views as authtoken_views
import data_ingest.urls

from validation.views import wifi_interceptor

urlpatterns = [
    path('admin/', admin.site.urls),
    url('accounts/', include('django.contrib.auth.urls')),

    url(r"^api-token-auth", authtoken_views.obtain_auth_token),
    url(r'^validate', wifi_interceptor, name="interceptor"),
]
