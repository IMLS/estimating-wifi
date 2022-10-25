from unittest import TestCase
import requests
import random
from datetime import datetime, timedelta

#testing variables
HOST = "http://localhost:3000"
NOW = datetime.now()
FSCS_ID = "AA0002-001"
LABEL = "Test in Room"


class FullTest ():
    #xkcd install key
    def GenKeyword(self):
        lines = open("5000-more-common.txt").read() 
        line = lines[0:] 
        words = line.split()
        mypass = ""
        for x in range(0, 4):
            myword = random.choice(words).upper()
            if x == 3:
                mypass += myword
            else:
                mypass += myword + "-" 
        #print passphrase
        return mypass

    #create JWT
    def test_jwt_gen(self, key, role):
        response = requests.post(
            f"{HOST}/rpc/jwt_gen", 
            json={
                "s_key" : key, 
                "s_role" : role
                }
        )
        #print(response.json())
        return response.json()

    #Post to sensor table
    def test_setup_post(self, FSCS_ID, LABEL, KEY, JWT):
        body = f'{{"_fscs":"{FSCS_ID}","_label": "{LABEL}", "_install_key": "{KEY}"}}'
        #print(body)
        #must be python dictionary
        headers = {
                    "content-type": "application/json", 
                    "Authorization": "Bearer " + JWT
                  }
        #print(headers)
        response = requests.post(
            f"{HOST}/rpc/sensor_setup", headers=headers, data=body
        )
        #print(response.json())
        return response.json()

    #Get sensor JWT token
    def test_info_post(self, SENSOR, KEY, JWT):
        body = f'{{"_sensor":{SENSOR}, "_install_key": "{KEY}"}}'
        #print(body)
        headers = {
                    "content-type": "application/json", 
                    "Authorization": "Bearer " + JWT
                  }
        response = requests.post(
            f"{HOST}/rpc/sensor_info", headers=headers, data=body
        )
        return response.json()

    #Post to heartbeat data
    def test_hb_post(self, SENSOR, FSCS, NOW, SERIAL, VERSION, JWT):
        body = f'{{"_sensor":{SENSOR},"_fscs":"{FSCS}","_hb": "{NOW}", "_serial": "{SERIAL}", "_version": "{VERSION}"}}'
        headers = {
                    "content-type": "application/json", 
                    "Authorization": "Bearer " + JWT
                  }
        response = requests.post(
            f"{HOST}/rpc/update_hb", headers=headers, data=body
        )
        return response.json()


#test setup workflow
#Test Secret
secret = ""
#Create class object
Test = FullTest()
#Create install key
key = Test.GenKeyword()
print(key)
#Create installer JWT
jwt = Test.test_jwt_gen(secret,"sensor")
jwt = jwt['token']
#Insert sensor data into table
setup = Test.test_setup_post(FSCS_ID, LABEL, key, jwt)
print(setup)
#Get sensor specific JWT and strip () from token
info = Test.test_info_post(setup, key, jwt)
info = info.strip("()")
print(info)
#Post a heartbeat using timestamp and new sensor_id
hb_test = Test.test_hb_post(setup, FSCS_ID, NOW, "1234ABC", "V.1A",info)
print(hb_test)

