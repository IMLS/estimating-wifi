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


def make_field(name, what, primary=False):
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


def create_table(table_name, fields, note=''):
    '''create a table with a primary key id'''
    data = {
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
    return requests.post(f'{URL}/collections/',
                         data=json.dumps(data),
                         headers=headers)


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

# delete if extant.
for table in ['wifi_raw', 'wifi_review', 'wifi_validated']:
    if table in tables:
        print(f'warning: deleting "{table}"')
        requests.delete(f'{URL}/collections/{table}', headers=headers)

print('building "wifi_raw"')
response = create_table(
    'wifi_raw',
    [
        make_field('date_created', 'timestamp'),
        make_field('data', 'json'),
        make_field('content_type', 'string'),
    ],
    note='raw wifi session data from rabbit',
)

print('building "wifi_review"')
response = create_table(
    'wifi_review',
    [
        make_field('date_created', 'timestamp'),
        make_field('headers', 'json'),
        make_field('whole_table_errors', 'json'),
        make_field('rows', 'json'),
        make_field('valid_row_count', 'integer'),
        make_field('invalid_row_count', 'integer'),
    ],
    note='wifi session data that did not pass rabbit validation',
)

# note: wifi_validated is deprecated
print('building "wifi_validated"')
response = create_table(
    'wifi_validated',
    [
        make_field('date_created', 'timestamp'),
        make_field('mac', 'string'),
        make_field('mfgs', 'string'),
        make_field('count', 'integer'),
    ],
    note='rabbit validated wifi session data',
)
