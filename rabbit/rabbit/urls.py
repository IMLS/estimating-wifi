from django.conf.urls import url
from django.urls import path

import data_ingest.urls  # required to avoid circular imports
from validation.views import wifi_interceptor

urlpatterns = [
    url(r'^validate/(?P<collection>[a-zA-Z0-9_]+)/$', wifi_interceptor, name="interceptor"),
]
