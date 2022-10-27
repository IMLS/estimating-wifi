from datetime import datetime, timedelta
import os
import requests
from unittest import TestCase


def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)


class PresencesTests(TestCase):
    def test_post_random_presence(self):
        token_url = endpoint(["rpc", "login"])
        body = {"fscs_id": "KY0069-002", "api_key": "hello-goodbye"}
        tr = requests.post(token_url, json=body)
        self.assertEqual(tr.status_code, 200)
        t0 = tr.json()["token"]

        now = datetime.now()
        url = endpoint(["rpc", "update_presence"])
        body = {
            "_fscs": "KY0069-002",
            "_start": str(now),
            "_end": str(now + timedelta(minutes=5)),
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json(), "KY0069-002")
