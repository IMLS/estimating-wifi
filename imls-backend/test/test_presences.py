from datetime import datetime, timedelta
import os
import requests
from unittest import TestCase
import urllib


def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)


class PresencesTests(TestCase):
    now = datetime.now()
    end = now + timedelta(minutes=5)

    def test_post_random_presence(self):
        token_url = endpoint(["rpc", "login"])
        body = {"fscs_id": "KY0069-002", "api_key": "hello-goodbye"}
        tr = requests.post(token_url, json=body)
        self.assertEqual(tr.status_code, 200)
        t0 = tr.json()["token"]

        url = endpoint(["rpc", "update_presence"])
        body = {
            "_start": str(self.now),
            "_end": str(self.end),
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json(), "KY0069-002")

    def test_verify_insertion(self):
        verify_insert_url = endpoint(["rpc", "verify_presence"])
        params = {"_fscs_id": "KY0069-002", "_start": str(self.now), "_end": str(self.end)}
        headers = {"Prefer": "count=estimated"}
        rv = requests.post(verify_insert_url, headers=headers, json=params)
        self.assertEqual(rv.status_code, 200)
        # FIXME: I'd like to know what the actual id is. There must be a better way
        # to validate that the insertion happened correctly. This query takes a 
        # particular library ID, start, and end time, and returns the UID for the presences table
        # if it exists. 
        print(rv.content)
        self.assertNotEqual(rv.content, b'null')
        self.assertIsInstance(int(rv.content), int)

    def test_verify_failure_on_non_insertion(self):
        verify_insert_url = endpoint(["rpc", "verify_presence"])
        params = {"_fscs_id": "KY0069-003", "_start": str(self.now), "_end": str(self.end)}
        headers = {"Prefer": "count=estimated"}
        rv = requests.post(verify_insert_url, headers=headers, json=params)
        self.assertEqual(rv.status_code, 200)
        print(rv.content)
        self.assertEqual(rv.content, b'null')
