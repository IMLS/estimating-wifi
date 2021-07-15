from django.conf.urls import url
from django.urls import path

import data_ingest.urls  # required to avoid circular imports

# monkey-patch for ReVal JSON ordering
from data_ingest import utils
def new_get_ordered_headers(headers):
    return sorted(headers)
utils.get_ordered_headers = new_get_ordered_headers

from validation.views import wifi_interceptor

urlpatterns = [
    url(r'^validate/(?P<collection>[a-zA-Z0-9_]+)/$', wifi_interceptor, name="interceptor"),
]
