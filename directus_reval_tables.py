import requests
import json
import sys
import time


URL = 'https://directus-demo2.app.cloud.gov'

INTERFACES = {
    'integer': 'numeric',
    'string': 'text-input',
    'timestamp': 'datetime',
    'json': None,
}


def jq(response):
    '''print in a format that `jq` can consume.'''
    print(json.dumps(response.json()))


def create_table_schema(table_name, fields, note=''):
    '''create a table with a primary key id'''
    return {
        'collection': table_name,
        'meta': {
            'icon' : 'build',
            'collection': table_name,
            'hidden': False,
            'note': note,
        },
        'fields': [
            {
                'field': 'id',
                'type': 'integer',
                'schema': {
                    'name': 'id',
                    'data_type': 'integer',
                    'is_primary_key': True,
                    'has_auto_increment': True,
                }
            }
        ] + fields,
    }


def create_field_schema(name, what, primary=False):
    if what not in ['string', 'integer', 'timestamp', 'json']:
        raise hell  # other types aren't supported here at present
    return {
        'collection': name,
        'field': name,
        'type': what,
        'display': 'formatted-json-value' if what == 'json' else None,
        'meta': {
            'collection': name,
            'field': name,
            'hidden': False,
            'interface': INTERFACES[what]
        }
    }


username = sys.argv[1]
password = sys.argv[2]

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

# delete.
for table in ['wifi_raw', 'wifi_review', 'wifi_validated']:
    if table in tables:
        print(f'warning: deleting "{table}"')
        requests.delete(f'{URL}/collections/{table}', headers=headers)

# start with wifi_raw.
print('building "wifi_raw"')
wifi_raw = create_table_schema(
    'wifi_raw',
    [
        create_field_schema('date_created', 'timestamp'),
        create_field_schema('data', 'json'),
        create_field_schema('content_type', 'string'),
    ],
    note='raw wifi session data from rabbit',
)
response = requests.post(f'{URL}/collections/', data=json.dumps(wifi_raw), headers=headers)
