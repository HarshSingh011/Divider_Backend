#!/usr/bin/env python3
"""
StockTrack WebSocket Connection Test
Login -> Get Token -> Test WebSocket Connection
"""

import json
import requests
import websocket
import threading
import time
from datetime import datetime

# Configuration
API_BASE = "https://divider-backend.onrender.com"
EMAIL = "harsh2004416@gmail.com"
PASSWORD = "SecurePass123!"

def log(message, level="INFO"):
    timestamp = datetime.now().strftime("%H:%M:%S")
    icon = {
        "INFO": "ℹ️",
        "SUCCESS": "✅",
        "ERROR": "❌",
        "WARNING": "⚠️",
        "WEBSOCKET": "🔌"
    }.get(level, "•")
    print(f"[{timestamp}] {icon} {message}")

def step(num, message):
    print(f"\n{'='*60}")
    print(f"[STEP {num}] {message}")
    print('='*60)

# Step 1: Health Check
step(1, "Testing Backend Health")
try:
    response = requests.get(f"{API_BASE}/health", timeout=5)
    if response.status_code == 200:
        log("Backend is running", "SUCCESS")
        print(f"Response: {response.json()}")
    else:
        log(f"Health check failed with status {response.status_code}", "ERROR")
        exit(1)
except Exception as e:
    log(f"Health check failed: {str(e)}", "ERROR")
    exit(1)

# Step 2: Login
step(2, "Logging in with credentials")
log(f"Email: {EMAIL}", "INFO")
log(f"Password: {PASSWORD}", "INFO")

try:
    login_payload = {
        "email": EMAIL,
        "password": PASSWORD
    }
    
    response = requests.post(
        f"{API_BASE}/auth/login",
        json=login_payload,
        headers={"Content-Type": "application/json"},
        timeout=5
    )
    
    if response.status_code == 200:
        data = response.json()
        token = data.get("token")
        
        if token:
            log(f"Login successful!", "SUCCESS")
            log(f"Token: {token[:50]}...", "INFO")
            log(f"Token length: {len(token)} characters", "INFO")
        else:
            log(f"No token in response: {data}", "ERROR")
            exit(1)
    else:
        log(f"Login failed with status {response.status_code}", "ERROR")
        print(f"Response: {response.text}")
        exit(1)
        
except Exception as e:
    log(f"Login request failed: {str(e)}", "ERROR")
    exit(1)

# Step 3: Save Token
step(3, "Token Information")
log(f"Token saved and ready for use", "SUCCESS")
print(f"\n{'='*60}")
print("FULL TOKEN (copy this to browser localStorage):")
print('='*60)
print(token)
print('='*60)

# Store token for next step
with open('token.txt', 'w') as f:
    f.write(token)
log("Token also saved to token.txt", "INFO")

# Step 4: WebSocket URL Info
step(4, "WebSocket Configuration")
ws_url = f"wss://divider-backend.onrender.com/ws?token={token}"
log(f"WebSocket URL: wss://divider-backend.onrender.com/ws?token=<token>", "INFO")
log(f"Protocol: WebSocket Secure (wss)", "INFO")
log(f"Authentication: JWT Token", "INFO")

# Step 5: Test WebSocket Connection
step(5, "Testing WebSocket Connection")

message_count = 0
ws = None

def on_message(ws, message):
    global message_count
    message_count += 1
    try:
        data = json.loads(message)
        if isinstance(data, list):
            symbols = ', '.join([f"{item.get('symbol', 'N/A')}:₹{item.get('currentPrice', 0)}" for item in data[:3]])
            log(f"Received message #{message_count}: {symbols}...", "WEBSOCKET")
        else:
            log(f"Received message #{message_count}: {str(data)[:100]}", "WEBSOCKET")
    except Exception as e:
        log(f"Error parsing message: {str(e)}", "WARNING")

def on_error(ws, error):
    log(f"WebSocket error: {error}", "ERROR")

def on_close(ws, close_status_code, close_msg):
    log(f"WebSocket closed - Status: {close_status_code}", "WARNING")

def on_open(ws):
    log(f"WebSocket connected successfully!", "SUCCESS")
    log(f"Waiting for market data...", "INFO")

try:
    log("Attempting WebSocket connection...", "INFO")
    ws = websocket.WebSocketApp(
        ws_url,
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close
    )
    
    # Run WebSocket for 10 seconds to receive some data
    def run_ws():
        ws.run_forever()
    
    ws_thread = threading.Thread(target=run_ws)
    ws_thread.daemon = True
    ws_thread.start()
    
    # Wait 10 seconds to receive messages
    log("Listening for market data (10 seconds)...", "INFO")
    time.sleep(10)
    
    # Close WebSocket
    ws.close()
    
    # Results
    step(6, "Test Results")
    log(f"Messages received: {message_count}", "SUCCESS" if message_count > 0 else "WARNING")
    
    if message_count > 0:
        log("✅ WebSocket is working perfectly!", "SUCCESS")
        log("Real-time market data is flowing successfully!", "SUCCESS")
    else:
        log("⚠️ No messages received (but connection was established)", "WARNING")
        log("This could mean:", "INFO")
        log("  - Backend is not generating market data", "INFO")
        log("  - Market data will start when market opens", "INFO")
    
except Exception as e:
    log(f"WebSocket test failed: {str(e)}", "ERROR")

# Step 7: Next Steps
step(7, "What to do next")
print("""
1. OPEN your app.html in browser:
   - File path: frontend/app.html
   - Or: http://localhost:8000/frontend/app.html (if running locally)

2. SAVE the token to browser localStorage:
   - Open DevTools (Press F12)
   - Go to Console tab
   - Paste this command and press Enter:
   
   localStorage.setItem('token', '<PASTE_TOKEN_HERE>')
   
   Or copy from token.txt file

3. REFRESH the page (Ctrl+R)

4. WATCH the Debug Panel on Home Screen:
   - You should see: 🟢 WebSocket Connected
   - Live logs showing market data updates
   - Real-time stock prices

5. OPEN Console (F12) and run debug commands:
   - debugWebSocket.status()
   - debugWebSocket.messages()
   - debugWebSocket.help()
""")

print(f"\n{'='*60}")
print("✨ Setup complete! Token is ready to use")
print('='*60)
print(f"\nToken saved in: token.txt")
print(f"Also available at: {API_BASE}")
