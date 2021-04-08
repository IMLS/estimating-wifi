from django.http import Http404, HttpResponseBadRequest
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from rest_framework import decorators, response
from rest_framework.parsers import JSONParser

import json
import requests
from data_ingest.ingestors import GoodtablesValidator

from rabbit import settings


def get_directus_validator(host, token, collection, version):
    url = f"https://{host}/items/validators_v{version}/{collection}"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}",
    }
    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Directus validation error: {response.content}")
    result = response.json()
    return result["data"]["validator"]


def proxy_data(host, token, collection, what, version):
    url = f"https://{host}/items/{collection}_v{version}"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}",
    }
    response = requests.post(url, data=json.dumps(what), headers=headers)
    if response.status_code != 200:
        raise Exception(f"Directus collection error: {response.content}")
    return response.json()


@csrf_exempt
@decorators.api_view(["POST"])
@decorators.parser_classes((JSONParser,))
def wifi_interceptor(request, collection=None):

    magic = request.META.get("HTTP_X_MAGIC_HEADER", None)
    if magic != settings.RABBIT_MAGIC_HEADER:
        raise Http404("Magic header not found")

    if not collection:
        raise Http404("Collection not found")

    if not request.data or not isinstance(request.data, list):
        return HttpResponseBadRequest("Data is malformed")

    host = request.META.get("HTTP_X_DIRECTUS_HOST", None)
    token = request.META.get("HTTP_X_DIRECTUS_TOKEN", None)
    if not host or not token:
        return HttpResponseBadRequest("Directus headers not found")

    # optional, defaults to 1
    version = request.META.get("HTTP_X_DIRECTUS_SCHEMA_VERSION", 1)

    # before we do anything else, save the raw data to directus.
    raw = {
        "collection": collection,
        "data": request.data,
        "content_type": request.content_type,
    }
    proxy_data(host, token, 'rabbit_raw', raw, version)

    # get the validator object.
    validator_json = get_directus_validator(host, token, collection, version)

    # bypass ReVal file checks.
    class InMemoryValidator(GoodtablesValidator):
        def load_file(self):
            return validator_json
    validator = InMemoryValidator('rabbit', 'temporary.csv')

    # TODO: SQL
    result = validator.validate(dict(source=request.data), request.content_type)

    # DESTINATION_FORMAT is not flexible enough in ReVal, so we proxy
    # the data manually -- we want to keep the raw data around and
    # either store the validated data or the data with errors.
    if result["valid"]:
        for item in request.data:
            proxy_data(host, token, collection, item, version)
    else:
        for table in result["tables"]:
            proxy_data(host, token, 'rabbit_review', table, version)

    return response.Response(result)
