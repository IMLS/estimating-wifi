from django.http import Http404
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from rest_framework import decorators, response
from rest_framework.parsers import JSONParser
from rest_framework.authentication import TokenAuthentication
from rest_framework.permissions import IsAuthenticated

import json
import requests
from data_ingest.ingestors import apply_validators_to
from . import settings


def get_directus_token():
    url = f"https://{settings.DIRECTUS_HOST}/auth/login"
    data = {"email": settings.DIRECTUS_USERNAME, "password": settings.DIRECTUS_PASSWORD}
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, data=json.dumps(data), headers=headers)
    result = response.json()
    if "errors" in result:
        raise Http404("Directus authentication error")
    return result["data"]["access_token"]


def proxy_data(token, collection, what):
    url = f"https://{settings.DIRECTUS_HOST}/items/{collection}"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}",
    }
    response = requests.post(url, data=json.dumps(what), headers=headers)
    result = response.json()
    if "errors" in result:
        raise Http404("Directus collection error")


@csrf_exempt
@decorators.api_view(["POST"])
@decorators.parser_classes((JSONParser,)) # TODO: support CSV
@decorators.authentication_classes([TokenAuthentication])
@decorators.permission_classes([IsAuthenticated])
def wifi_interceptor(request):
    # before we do anything else, save the raw data to directus.
    token = get_directus_token()
    proxy_data(token, 'wifi_raw', {
        "data": request.data,
        "content_type": request.content_type
    })

    result = apply_validators_to(request.data, request.content_type)
    # DESTINATION_FORMAT is not flexible enough in ReVal, so we proxy
    # the data manually -- we want to keep the raw data around and
    # either store the validated data or the data with errors.
    if result["valid"]:
        for item in request.data["source"]:
            proxy_data(token, 'wifi_validated', item)
    else:
        for table in result["tables"]:
            proxy_data(token, 'wifi_review', table)

    return response.Response(result)
