import requests
import json
import sys

username = sys.argv[1]
password = sys.argv[2]

base_url = "http://localhost:8000"

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

content = json.dumps({"source": [
    {"mac": "60:38:e0", "mfgs": "Belkin", "count": 20},
    {"mac": "something", "mfgs": "unknown", "count": 20},
]})
resp = requests.post(f"{base_url}/validate/", data=content, headers=headers)
validation = resp.json()

print([row["errors"][0]["message"] for row in validation["tables"][0]['rows']])
