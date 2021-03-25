import requests
import json
import sys
import csv

username = sys.argv[1]
password = sys.argv[2]
filename = sys.argv[3]

base_url = "https://10x-rabbit-demo.app.cloud.gov"

resp = requests.post(
    f"{base_url}/api-token-auth/", data=dict(username=username, password=password),  # nosec
)

token_result = resp.json()
if "token" not in token_result:
    print(f"error: could not obtain token: {token_result}")
    sys.exit(1)
token = token_result["token"]
print(f"obtained token: {token}")

headers = {
    "Content-Type": "application/json",
    "Authorization": f"Token {token}",
}

data = []
csv_file = csv.reader(open(filename, 'r'))
csv_headers = next(csv_file)
result = [dict(zip(csv_headers, line)) for line in csv_file]

content = json.dumps({"source": result})

def print_error(resp):
    tables = resp['tables']
    for table in tables:
        for row in table["rows"]:
            errors = row["errors"]
            if errors:
                messages = ';'.join([e["message"] for e in errors])
                print(f'Row number {row["row_number"]}: {messages}')

resp = requests.post(f"{base_url}/validate/", data=content, headers=headers)
validation = resp.json()
if validation["valid"]:
    print("CSV file passed validation")
else:
    print("CSV file failed with the following errors:")
    print_error(validation)
