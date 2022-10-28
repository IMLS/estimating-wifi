import os
import requests
from unittest import TestCase

def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)

class IMLSTests(TestCase):
    
    def test_actual_library_timezone(self):
        url = endpoint(["rpc", "get_library_timezone"])
        query = {"fscs_id": "KY0069-002"}
        r = requests.post(url, json=query)
        print(r.json())
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['time'], "00:00:00-04")

    def test_timezone_lookup_fail(self):
        url = endpoint(["rpc", "get_library_timezone"])
        query = {"fscs_id": "KY0069-003"}
        r = requests.post(url, json=query)
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['time'], None)

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
        self.assertGreater(estimated_count, 100)
