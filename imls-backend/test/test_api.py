from unittest import TestCase
from datetime import datetime, timedelta
import requests
import pytz



HOST = "http://localhost:3000"
NOW = datetime.now()
NOW_5 = NOW + timedelta(minutes=5)
EX_SENSOR_ID = 10
FSCS_ID = "AA0001-001"
EX_INSTALL_KEY = "blue-red-dog-bird"
NEW_INSTALL_KEY = "orange-black-xxx-yyy"
JWT = ""


class Presences(TestCase):
    def test_presences(self):
        response = requests.get(f"{HOST}/presences")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)

    def test_post_presences(self):
        body = f'{{"_start": "{NOW}", "_end": "{NOW_5}", \
                   "_fscs":"{FSCS_ID}","_sensor":{EX_SENSOR_ID}, \
                   "_serial": "123A", "_version": "X999"}}'
        token = f'{{"Authorization": "{JWT}"}}'
        response = requests.post(f"{HOST}/update_presence", headers=token, data=body)
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)


class BinDevices(TestCase):
    def test_day_query(self):
        response = requests.get(
            f"{HOST}/rpc/bin_devices_per_hour?_start=2022-05-10&_fscs_id={FSCS_ID}"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 24)

    def test_multi_day_query(self):
        response = requests.get(
            f"{HOST}/rpc/bin_devices_over_time?_start=2022-05-10&_fscs_id={FSCS_ID}&_direction=true&_days=7"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 24)

class Heartbeats(TestCase):
    def test_heartbeats(self):
        response = requests.get(f"{HOST}/presences")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)

    def test_hb_post(self):
        body = f'{{"_sensor":{EX_SENSOR_ID},"_fscs":"{FSCS_ID}","_hb": "{NOW}", "_serial": "123A", "_version": "X999"}}'
        token = f'{{"Authorization": "{JWT}"}}'
        response = requests.post(
            f"{HOST}/rpc/update_hb", headers=token, data=body
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)

class JWT(TestCase):
    def test_jwt_gen(self):
        body = '{"s_key": "BlahBlahBlahBlahBlahBlahBlahBlah, "s_role" : "sensor"}'
        response = requests.post(
            f"{HOST}/rpc/jwt_gen", data=body
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)

class LibSearch(TestCase):
    def test_imls_query(self):
        response = requests.get(f"{HOST}/imls_lookup")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)
    
    def test_fscs_query(self):
        response = requests.get(
            f"{HOST}/rpc/lib_search_fscs?_fscs_id={FSCS_ID}"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)

    def test_state_query(self):
        response = requests.get(
            f"{HOST}/rpc/lib_search_state?_state_code=AK"
        )
        self.assertTrue(response.status_code > 0)
        items = response.json()
        self.assertTrue(len(items) > 0)
    
    def test_name_query(self):
        response = requests.get(
            f"{HOST}/rpc/lib_search_name?_name=POINT"
        )
        self.assertTrue(response.status_code > 0)
        items = response.json()
        self.assertTrue(len(items) > 0)

class Sensor(TestCase):
    def test_setup_post(self):
        body = f'{{"_fscs":"{FSCS_ID}","_label": "test label", "_install_key": "{NEW_INSTALL_KEY}"}}'
        token = f'{{"Authorization": "{JWT}"}}'
        response = requests.post(
            f"{HOST}/rpc/sensor_setup", headers=token, data=body
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)

    def test_info_post(self):
        body = f'{{"_sensor":"{EX_SENSOR_ID}", "_install_key": "{EX_INSTALL_KEY}"}}'
        token = f'{{"Authorization": "{JWT}"}}'
        response = requests.post(
            f"{HOST}/rpc/bin_devices_per_hour?_start=2022-05-10&_fscs_id=AA0003-001"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)
