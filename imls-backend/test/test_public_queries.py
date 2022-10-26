import os
import requests
from unittest import TestCase

scheme = "http"
hostname = "localhost"
port = 3000
test_url = f"{scheme}://{hostname}:{port}"

def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)

class IMLSTests(TestCase):
    def test_imls_query(self):
        url = endpoint(["imls_lookup"])
        print("URL: ", url)
        response = requests.get(url)
        if response.status_code != 200:
            print(response.json())
        self.assertTrue(response.status_code == 200)
        items = response.json()
        if len(items) <= 0:
            print("ITEMS: ", items)
        self.assertTrue(len(items) > 0)
    
