from django.http import Http404, HttpResponseBadRequest
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from rest_framework import decorators, response
from rest_framework.parsers import JSONParser

import json
import requests
from data_ingest.ingestors import apply_validators_to

from rabbit import settings


def get_directus_token(host, email, password):
    url = f"https://{host}/auth/login"
    data = {"email": email, "password": password}
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, data=json.dumps(data), headers=headers)
    result = response.json()
    if "errors" in result:
        return None
    return result["data"]["access_token"]


def get_directus_validator(host, token, collection):
    url = f"https://{host}/items/validators/{collection}"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}",
    }
    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Directus validation error: {response.content}")
    return response.json()


def proxy_data(host, token, collection, what):
    url = f"https://{host}/items/{collection}"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}",
    }
    response = requests.post(url, data=json.dumps(what), headers=headers)
    if response.status_code != 200:
        raise Exception(f"Directus collection error: {response.content}")


@csrf_exempt
@decorators.api_view(["POST"])
@decorators.parser_classes((JSONParser,)) # TODO: support CSV
def wifi_interceptor(request, collection=None):

    magic = request.META.get("HTTP_X_MAGIC_HEADER", None)
    if magic != settings.MAGIC_HEADER:
        raise Http404("Magic header not found")

    if not collection:
        raise Http404("Collection not found")

    if not request.data or not isinstance(request.data, list):
        return HttpResponseBadRequest("Data is malformed")

    host = request.META.get("HTTP_X_DIRECTUS_HOST", None)
    email = request.META.get("HTTP_X_DIRECTUS_EMAIL", None)
    password = request.META.get("HTTP_X_DIRECTUS_PASSWORD", None)
    if not host or not email or not password:
        return HttpResponseBadRequest("Directus headers not found")

    token = get_directus_token(host, email, password)
    if not token:
        return HttpResponseBadRequest("Directus authentication error")

    # TODO: get tests working
    # TODO: grab json validation and use that instead of wifi.json

    # before we do anything else, save the raw data to directus.
    raw = {
        "collection": collection,
        "data": request.data,
        "content_type": request.content_type,
        # TODO: validation
    }
    proxy_data(host, token, 'rabbit_raw', raw)

    result = apply_validators_to(dict(source=request.data), request.content_type)

    # DESTINATION_FORMAT is not flexible enough in ReVal, so we proxy
    # the data manually -- we want to keep the raw data around and
    # either store the validated data or the data with errors.
    if result["valid"]:
        for item in request.data:
            proxy_data(host, token, collection, item)
    else:
        for table in result["tables"]:
            proxy_data(host, token, 'rabbit_review', table)

    return response.Response(result)
