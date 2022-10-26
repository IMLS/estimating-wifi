from datetime import datetime, timedelta
import jwt
import requests
import pytz
from unittest import TestCase



HOST = "http://localhost:3000"
NOW = datetime.now()
NOW_5 = NOW + timedelta(minutes=5)
EX_SENSOR_ID = 10
FSCS_ID = "AA0005-001"
EX_INSTALL_KEY = "blue-red-dog-bird"
NEW_INSTALL_KEY = "orange-black-xxx-yyy"
JWTSECRET = "DozeDischargeLadderStriveUnthawedCharting"

def generate_jwt(role):
    payload = {"role": str(role), "email": "anyone@anywhere.com"}
    token = jwt.encode(payload=payload, key=JWTSECRET,  algorithm="HS256")
    print(f"TOKEN: |{token}|")
    return token

class PresencesTests(TestCase):
    def test_presences(self):
        response = requests.get(f"{HOST}/presences")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) > 0)
        print("test_presences PASSED")

    # This function is currently marked as a 'postgres' function.
    # Does that mean we cannot call it via the API externally?
    # If so, we need to test this differently. Or, we need to create a user
    # that we authenticate as for *testing*, and test it that way.
    # def test_post_presences(self):
    #     body = {
    #         "_start": str(NOW), 
    #         "_end": str(NOW_5), 
    #         "_fscs": str(FSCS_ID),
    #         "_sensor": str(EX_SENSOR_ID),
    #         "_serial": "123A", 
    #         "_version": "X999"
    #         }
    #     token = { "Authorization": f"Bearer {generate_jwt('sensor')}"}
    #     response = requests.post(f"{HOST}/rpc/update_presence", headers=token, json=body)
    #     print(response.json())
    #     self.assertTrue(response.status_code == 200)
    #     items = response.json()
    #     self.assertTrue(len(items) == 1)


class BinDevicesTests(TestCase):
    def test_day_query(self):
        response = requests.get(
            f"{HOST}/rpc/bin_devices_per_hour?_start=2022-05-10&_fscs_id={FSCS_ID}"
        )
        self.assertTrue(response.status_code == 200)
        items = response.json()
        print(f"ITEMS: {items}")
        self.assertTrue(len(items) == 24)

    def test_multi_day_query(self):
        _start = "2022-05-10"
        _fscs_id = "AA0003-001"
        _direction = "true"
        _days = "7"
        url = (f"{HOST}/rpc/bin_devices_over_time" + 
            f"?_start={_start}" +
            f"&_fscs_id={_fscs_id}" + 
            f"&_direction={_direction}" + 
            f"&_days={_days}")
        print(f"URL: |{url}|")
        response = requests.get(url)
        print(f"RESPONSE: |{response}|")
        self.assertTrue(response.status_code == 200)
        items = response.json()
        print(f"ITEMS: {items}")
        self.assertEquals(list(map(len, items)), [24, 24, 24, 24, 24, 24, 24])

class HeartbeatTests(TestCase):
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

class JWTTests(TestCase):
    def test_jwt_gen(self):
        response = requests.post(
            f"{HOST}/rpc/jwt_gen", 
            json={
                "s_key" : "BlahBlahBlahBlahBlahBlahBlahBlah", 
                "s_role" : "sensor"
                }
        )
        print(response.json())
        self.assertTrue(response.status_code == 200)
        items = response.json()
        self.assertTrue(len(items) == 1)

class LibSearchTests(TestCase):
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

class SensorTests(TestCase):
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
