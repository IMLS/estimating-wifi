import os
import requests
from unittest import TestCase

def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)

class IMLSTests(TestCase):
    def test_existence_of_libraries_in_imls_lookup_table(self):
        url = endpoint(["imls_lookup"])
        response = requests.get(url)
        # If we don't see a 200 response, that's just plain bad.
        if response.status_code != 200:
            print(response.json())
        self.assertTrue(response.status_code == 200)
        items = response.json()
        # We should always see multiple libraries here, even in production.
        # If we don't that means something is very broken.
        self.assertTrue(len(items) > 0)
    
    # NOTE: There's a ?limit parameter on these, because otherwise a lot of 
    # values come back by default, and that makes for  slow tests.
    def test_existence_of_presences(self):
        url = endpoint(["presences?limit=25"])
        response = requests.get(url)
        if response.status_code != 200:
            print(response.json())
        self.assertTrue(response.status_code == 200)
    
    def test_presences_is_big(self):
        url = endpoint(["presences?limit=25"])
        headers = {"Prefer": "count=estimated"}
        response = requests.get(url, headers=headers)
        estimated_count = int(response.headers["Content-Range"].split("/")[1])
        self.assertGreater(estimated_count, 10000)
