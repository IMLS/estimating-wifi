import datetime
import os
import pytz
import requests
from unittest import TestCase

def endpoint(ep_arr):
    test_url = f"{os.getenv('POETRY_SCHEME')}://{os.getenv('POETRY_HOSTNAME')}:{os.getenv('POETRY_PORT')}"
    return test_url + "/" + "/".join(ep_arr)

class IMLSTests(TestCase):
    def test_existence_of_libraries_in_imls_lookup_table(self):
        url = endpoint(["imls_lookup"])
        response = requests.get(url)
        # If we don't see a 200 response, that's just plain bad.
        if response.status_code != 200:
            print(response.json())
        self.assertTrue(response.status_code == 200)
        items = response.json()
        # We should always see multiple libraries here, even in production.
        # If we don't that means something is very broken.
        self.assertTrue(len(items) > 0)

class PresencesTests(TestCase):
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


def validate(test_obj, fscs_id):
    url = endpoint(["rpc", "lib_search_fscs"])
    query = {"_fscs_id": fscs_id}
    r = requests.post(url, json=query)
    obj = r.json()
    test_obj.assertEqual(obj['fscskey'], fscs_id.split("-")[0])
class LibSearchTests(TestCase):
    def test_validate_one(self):
        url = endpoint(["rpc", "lib_search_fscs"])
        query = {"_fscs_id": "KY0069-002"}
        r = requests.post(url, json=query)
        print(r.json())
        obj = r.json()
        self.assertEqual(obj['fscskey'], "KY0069")

    # First, this takes forever.
    # Second, I generated the list of FSCS Ids from the DB itself.
    # That means they're all "good." This loop tests that my list
    # came from the DB, and is still there. Not useful.
    # def test_validate_all(self):
    #     for line in open("all_fscs_ids.txt"):
    #         self.assertEqual(len(line.strip()), 10)
    #         validate(self, line.strip())

class BinningTests(TestCase):
    def test_no_data(self):
        # _start date, _fscs_id text, _direction boolean, _days integer
        url = endpoint(["rpc", "bin_devices_over_time"])
        query = {
            "_fscs_id": "KY0069-002",
            "_start": "2022-02-02",
            "_direction": True,
            "_days": 7    
                }
        r = requests.post(url, json=query)
        print(r.json())
        res = r.json()
        # This should return 7 lists of length 24, each containing all zeros.
        self.assertEqual(7, len(res))
        map(lambda ls: self.assertEqual(len(ls), 24), res)
        map(lambda ls: map(lambda v: self.assertEqual(v, 0), ls), res)

    def test_actual_over_time(self):
        url = endpoint(["rpc", "bin_devices_over_time"])
        query = {
            "_fscs_id": "AA0005-001",
            "_start": "2022-05-11",
            "_direction": True,
            "_days": 7    
                }
        r = requests.post(url, json=query)
        print(r.json())
        res = r.json()
        # This should return 7 lists of length 24, each containing all zeros.
        self.assertEqual(7, len(res))
        actual = [
            [156, 156, 156, 150, 144, 150, 150, 144, 138, 162, 126, 120, 108, 114, 138, 102, 120, 90, 108, 120, 108, 90, 90, 84],
            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 90, 108, 108, 102, 102], 
            [114, 126, 120, 144, 138, 108, 96, 102, 114, 126, 126, 132, 108, 120, 114, 126, 108, 90, 90, 96, 78, 78, 78, 72], 
            [72, 78, 78, 78, 84, 84, 90, 90, 78, 84, 84, 96, 132, 90, 72, 84, 78, 60, 66, 66, 60, 54, 48, 48], 
            [48, 54, 48, 48, 48, 54, 54, 60, 60, 72, 66, 66, 60, 72, 72, 66, 84, 84, 66, 60, 60, 48, 48, 48], 
            [0, 0, 0, 0, 36, 66, 60, 66, 60, 60, 78, 78, 60, 60, 54, 48, 48, 54, 60, 48, 48, 48, 60, 36], 
            [60, 60, 54, 54, 54, 54, 54, 54, 42, 48, 54, 66, 66, 60, 60, 66, 78, 96, 96, 96, 66, 72, 66, 36]
            ]
        self.assertEqual(res, actual)

    def test_actual_per_hour(self):
        url = endpoint(["rpc", "bin_devices_per_hour"])
        query = {
            "_fscs_id": "AA0005-001",
            "_start": "2022-05-11"
                }
        r = requests.post(url, json=query)
        print(r.json())
        res = r.json()
        # This should return 7 lists of length 24, each containing all zeros.
        self.assertEqual(24, len(res))
        actual = [156, 156, 156, 150, 144, 150, 150, 144, 138, 162, 126, 120, 108, 114, 138, 102, 120, 90, 108, 120, 108, 90, 90, 84]
        self.assertEqual(res, actual)