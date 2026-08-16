import requests

url = "http://db-api:8080/"

try:
    r = requests.get(url, timeout=5)
    print(f"Connected! HTTP {r.status_code}")
    print(r.text)
except requests.RequestException as e:
    print(f"Failed to connect: {e}")