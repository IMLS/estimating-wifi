import requests
import json
import sys
import time


URL = 'https://directus-demo2.app.cloud.gov'

RESERVED_TABLES = [
    'pls_data',
    'pls_events',
    'wifi_raw',
    'wifi_review',
    'wifi_validated',
]

INTERFACES = {
    'integer': 'numeric',
    'string': 'text-input',
    'timestamp': 'datetime',
}


def jq(response):
    '''print in a format that `jq` can consume.'''
    print(json.dumps(response.json()))


def create_field(name, what, primary=False):
    if what not in ['string', 'integer', 'timestamp']:
        raise hell  # other types aren't supported here at present
    return {
        'collection': table_name,
        'field': name,
        'type': what,
        'meta': {
            'collection': table_name,
            'field': name,
            'hidden': False,
            'interface': INTERFACES[what]
        }
    }


username = sys.argv[1]
password = sys.argv[2]
table_name = sys.argv[3] if len(sys.argv) > 3 else 'will_it_blend'
table_name = table_name.lower()

if table_name in RESERVED_TABLES:
    print(f'error: "{table_name}" is a reserved table.')
    sys.exit(1)

# authenticate.
data = {'email': username, 'password': password}
headers = {'Content-Type': 'application/json'}
response = requests.post(f'{URL}/auth/login', data=json.dumps(data), headers=headers)
respjson = response.json()
token = respjson['data']['access_token']
headers['Authorization'] = f'Bearer {token}'

# query for extant tables.
response = requests.get(f'{URL}/collections', headers=headers)
tables = [c['collection'] for c in response.json()['data']]

# delete if requested.
if len(sys.argv) == 5 and sys.argv[4] == 'delete':
    if table_name in tables:
        print(f'deleting "{table_name}"')
        requests.delete(f'{URL}/collections/{table_name}', headers=headers)

# otherwise, create.
if table_name not in tables:
    magic_fields = [
        {
            'field': 'magic_index',
            'type': 'integer',
            'schema': {
                'is_primary_key': True,
                'has_auto_increment': True
            }
        }
    ]
    data = {
        'collection': table_name,
        'meta': {
            'icon' : 'check_circle',
            'collection': table_name,
            'hidden': False
        },
        'fields': magic_fields
    }
    response = requests.post(f'{URL}/collections/', data=json.dumps(data), headers=headers)
    print(jq(response))

# add fields.
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

for field in pls_data_fields:
    response = requests.post(f'{URL}/fields/{table_name}', data=json.dumps(field), headers=headers)

# response = requests.get(f'{URL}/fields/people2', headers=headers)
# print(json.dumps(response.json()))
