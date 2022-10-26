import jwt
import requests
from unittest import TestCase

scheme = "http"
hostname = "localhost"
port = 3000
test_url = f"{scheme}://{hostname}:{port}"

JWTSECRET = "DozeDischargeLadderStriveUnthawedCharting"

def endpoint(ep_arr):
    return test_url + "/" + "/".join(ep_arr)

def generate_jwt(role):
    payload = {"role": str(role), "email": "anyone@anywhere.com"}
    token = jwt.encode(payload=payload, key=JWTSECRET,  algorithm="HS256")
    print(f"TOKEN: |{token}|")
    return token

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
    
