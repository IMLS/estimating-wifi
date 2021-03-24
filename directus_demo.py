import requests
import json
import sys
import time


URL = 'https://directus-demo.app.cloud.gov'
RESERVED_TABLES = [
    'people',
    'people2',
    'wifi_raw',
    'wifi_review',
    'wifi_validated',
    'program_attendance'
]


def jq(response):
    """print in a format that `jq` can consume."""
    print(json.dumps(response.json()))


username = sys.argv[1]
password = sys.argv[2]
table_name = sys.argv[3] if len(sys.argv) > 3 else 'test_table'
table_name = table_name.lower()

if table_name in RESERVED_TABLES:
    print(f'error: "{table_name}" is a reserved table.')
    sys.exit(1)

# authenticate.
data = {'email': username, 'password': password}
headers = {'Content-Type': 'application/json'}
response = requests.post(f'{URL}/auth/login', data=json.dumps(data), headers=headers)

token = response.json()['data']['access_token']
headers['Authorization'] = f'Bearer {token}'

# create table anew.
response = requests.get(f'{URL}/collections', headers=headers)
tables = [c['collection'] for c in response.json()['data']]

# # deleting collections doesn't work properly in directus.
# if table_name in tables:
#     print(f'deleting "{table_name}"')
#     requests.delete(f'{URL}/collections/{table_name}', headers=headers)

if table_name not in tables:

    def create_field(name, what, **kwargs):
        if what not in ['string', 'integer', 'timestamp']:
            raise hell  # other types aren't supported here at present
        return {"field": name, "type": what, **kwargs}

    pls_data_fields = [
        create_field('event_id', 'integer'),
        create_field('device_uuid', 'string'),
        create_field('lib_user', 'string'),
        create_field('localtime', 'timestamp'),
        create_field('servertime', 'timestamp'),
        create_field('session_id', 'integer'),
        create_field('device_id', 'string'),
        create_field('last_seen', 'integer'),
    ]
    data = {"collection": table_name, "fields": pls_data_fields}
    response = requests.post(f'{URL}/collections/', data=json.dumps(data), headers=headers)
    print(jq(response))

# response = requests.get(f'{URL}/fields/people2', headers=headers)
# print(json.dumps(response.json()))
