from unittest import TestCase
import requests


HOST = "http://localhost:3000"


class Presences(TestCase):
    def test_presences(self):
        response = requests.get(f"{HOST}/presences")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)


class BinDevices(TestCase):
    def test_day_query(self):
        response = requests.get(
            f"{HOST}/rpc/bin_devices_per_hour?_start=2022-05-10&_fscs_id=AA0003-001"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 24)
