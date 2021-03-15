from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt
from rest_framework import decorators, response
from rest_framework.parsers import JSONParser
from rest_framework.authentication import TokenAuthentication
from rest_framework.permissions import IsAuthenticated

from data_ingest.ingestors import apply_validators_to


@csrf_exempt
@decorators.api_view(["POST"])
@decorators.parser_classes((JSONParser,)) # TODO: support CSV
@decorators.authentication_classes([TokenAuthentication])
@decorators.permission_classes([IsAuthenticated])
def wifi_interceptor(request):
    result = apply_validators_to(request.data, request.content_type)
    # TODO: pass on to directus
    return response.Response(result)
