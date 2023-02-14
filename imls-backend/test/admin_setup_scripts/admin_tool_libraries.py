import random
import requests
import admin_tool_kit_checks as ck


def callSetup(s, u, p, fscs_id):
    """
    Call "sensor_setup" to save info to the database.
    """
    token_url = s + "/rpc/login"
    body = {"fscs_id": u, "api_key": p}
    # Need to post, not get, if you're passing params in the body.
    tr = requests.post(token_url, json=body)
    t0 = tr.json()["token"]
    url = s + "/rpc/library_setup" 
    body = {
        "_fscs_id": fscs_id,
    }
    headers = {"Authorization": f"Bearer {t0}"}
    r = requests.post(url, json=body, headers=headers)
    print("UPSERT RESPONSE", r.json())

def callRemoveLibrary(s, u, p, fscs_id):
    """
    Call "sensor_remove" to remove a sensor from the database.
    """
    token_url = s + "/rpc/login"
    body = {"fscs_id": u, "api_key": p}
    # Need to post, not get, if you're passing params in the body.
    tr = requests.post(token_url, json=body)
    t0 = tr.json()["token"]
    url = s + "/rpc/library_remove" 
    body = {
        "_fscs_id": fscs_id,
    }
    headers = {"Authorization": f"Bearer {t0}"}
    r = requests.post(url, json=body, headers=headers)
    print("SERVER RESPONSE", r.json())