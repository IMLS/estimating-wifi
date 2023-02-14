import requests

def callSetup(s, u, p, sensor, key, label):
    """
    Call "sensor_setup" to save info to the database.
    """
    token_url = s + "/rpc/login"
    body = {"fscs_id": u, "api_key": p}
    # Need to post, not get, if you're passing params in the body.
    tr = requests.post(token_url, json=body)
    t0 = tr.json()["token"]
    url = s + "/rpc/sensor_setup" 
    body = {
        "_sensor": sensor,
        "_key": key,
        "_label": label,
    }
    headers = {"Authorization": f"Bearer {t0}"}
    r = requests.post(url, json=body, headers=headers)
    print("UPSERT RESPONSE", r.json())

def callUpdatePassword(s, u, p, sensor, key):
        """
        Call "update_password" to update sensor password in the database.
        """
        token_url = s + "/rpc/login"
        body = {"fscs_id": u, "api_key": p}
        # Need to post, not get, if you're passing params in the body.
        tr = requests.post(token_url, json=body)
        t0 = tr.json()["token"]
        url = s + "/rpc/update_password" 
        body = {
            "_sensor": sensor,
            "_key": key,
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        print("UPDATE RESPONSE", r.json())

def callRemoveSensor(s, u, p, sensor):
        """
        Call "sensor_remove" to remove a sensor from the database.
        """
        token_url = s + "/rpc/login"
        body = {"fscs_id": u, "api_key": p}
        # Need to post, not get, if you're passing params in the body.
        tr = requests.post(token_url, json=body)
        t0 = tr.json()["token"]
        url = s + "/rpc/sensor_remove" 
        body = {
            "_sensor": sensor,
        }
        headers = {"Authorization": f"Bearer {t0}"}
        r = requests.post(url, json=body, headers=headers)
        print("SERVER RESPONSE", r.json())

