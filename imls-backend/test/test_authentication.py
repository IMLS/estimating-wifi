import os
import requests
from unittest import TestCase


def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)


class LoginTests(TestCase):
    def test_call_login_invalid(self):
        url = endpoint(["rpc", "login"])
        body = {"fscs_id": "ME0018-001", "api_key": "notapassword"}
        # Need to post, not get, if you're passing params in the body.
        r = requests.post(url, json=body)
        print(r.json())
        self.assertEqual(r.status_code, 403)
        # Should be invalid user or password error
        self.assertEqual(r.json()["code"], "28P01")

    def test_login_successfully(self):
        url = endpoint(["rpc", "login"])
        body = {"fscs_id": "KY0069-002", "api_key": "hello-goodbye"}
        # Need to post, not get, if you're passing params in the body.
        r = requests.post(url, json=body)
        print("URL ", url)
        print("RESPONSE ", r.json())
        self.assertEqual(r.status_code, 200)
        t0 = r.json()["token"]
        # This comes back as three parts separated by periods.
        # blasdlkfjasdflkjas.balskdjfaskjdfhaksjdhf.baskdjhfaskdjhfashkldj
        # The first part of the token will always be the same.
        t1 = t0.split(".")[0]
        _ = t0.split(".")[1]
        self.assertEqual(t1, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")

    def test_use_token(self):
        """
        Call "beat_the_heart" to see if we can get through authentication.
        """
        token_url = endpoint(["rpc", "login"])
        body = {"fscs_id": "KY0069-002", "api_key": "hello-goodbye"}
        # Need to post, not get, if you're passing params in the body.
        tr = requests.post(token_url, json=body)
        print("URL ", token_url)
        print("RESPONSE ", tr.json())
        self.assertEqual(tr.status_code, 200)
        t0 = tr.json()["token"]

        url = endpoint(["rpc", "beat_the_heart"])
        body = {
            "_sensor_serial": "abcde",
            "_sensor_version": "3.0",
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        print(r.json())
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()["result"], "OK")
