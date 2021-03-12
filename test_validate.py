import requests
import json
import sys

username = sys.argv[1]
password = sys.argv[2]

base_url = "http://localhost:8000/data/api"

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

content = { "source": [] }
resp = requests.post(f"{base_url}/validate/", json=content, headers=headers)
print(resp)
print(resp.json())
