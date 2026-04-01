import requests
import json

try:
    r = requests.post('https://divider-backend.onrender.com/auth/login', 
        json={'email':'harsh2004416@gmail.com','password':'SecurePass123!'},
        timeout=10)
    data = r.json()
    token = data.get('token', '')
    
    with open('c:/Users/harsh/Python/StockTrack/token_output.txt', 'w') as f:
        f.write(f"Status: {r.status_code}\n")
        f.write(f"Token: {token}\n")
        f.write(f"Response: {json.dumps(data, indent=2)}\n")
except Exception as e:
    with open('c:/Users/harsh/Python/StockTrack/token_output.txt', 'w') as f:
        f.write(f"Error: {str(e)}\n")
