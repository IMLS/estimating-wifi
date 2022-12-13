from datetime import datetime, timedelta
import os
import requests
from unittest import TestCase
import urllib


def endpoint(ep_arr):
    test_url = f"{os.getenv('PYTEST_SCHEME')}://{os.getenv('PYTEST_HOSTNAME')}:{os.getenv('PYTEST_PORT')}"
    return test_url + "/" + "/".join(ep_arr)


class TimezoneTests(TestCase):
    def test_get_timezone_from_fscs_id(self):
        url = endpoint(["rpc", "get_timezone_from_fscs_id"])
        params = {"_fscs_id": "KY0069-002"}
        response = requests.post(url, json=params)
        # If we don't see a 200 response, that's just plain bad.
        print(response.json())
        self.assertTrue(response.status_code == 200)
        tz = response.json()
        # We should always see multiple libraries here, even in production.
        # If we don't that means something is very broken.
        self.assertEqual(tz, -4)

    def test_tz_fail_on_non_id(self):
        url = endpoint(["rpc", "get_timezone_from_fscs_id"])
        params = {"_fscs_id": "YK0000-002"}
        response = requests.post(url, json=params)
        # If we don't see a 200 response, that's just plain bad.
        print(response.json())
        self.assertTrue(response.status_code == 200)
        tz = response.json()
        # We should always see multiple libraries here, even in production.
        # If we don't that means something is very broken.
        self.assertEqual(tz, None)
