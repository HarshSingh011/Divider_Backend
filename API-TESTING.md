# StockTrack API - Quick Testing Guide

## Quick Start (Copy & Paste Commands)

### 1. Check Server Health
```bash
curl -X GET http://localhost:8080/health
```

### 2. Register New User
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"harsh@example.com",
    "username":"harsh",
    "password":"password123"
  }'
```

**Response:**
```json
{
  "id": "user123",
  "email": "harsh@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Login User (Get Token)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"harsh@example.com",
    "password":"password123"
  }'
```

**Save the token:** Copy the token value from response for use in protected endpoints

---

## Protected Endpoints (Require Token)

### 4. Get Wallet Snapshot
```bash
curl -X GET http://localhost:8080/trading/wallet \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Response:**
```json
{
  "userId": "user123",
  "totalBalance": 99698.75,
  "investedAmount": 10000,
  "cashAvailable": 89698.75,
  "positions": [
    {
      "symbol": "RELIANCE",
      "quantity": 100,
      "averagePrice": 100.50,
      "currentPrice": 102.50,
      "profitLoss": 200.00
    }
  ],
  "recentTransactions": [...]
}
```

### 5. Get OHLC Candlesticks (for Charts)
```bash
curl -X GET "http://localhost:8080/trading/candles?symbol=RELIANCE&limit=50" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Response:**
```json
[
  {
    "symbol": "RELIANCE",
    "timestamp": "2026-03-30T14:40:00Z",
    "open": 102.00,
    "high": 103.50,
    "low": 101.50,
    "close": 103.20,
    "volume": 1523400
  },
  ...
]
```

### 6. Execute Trade (Buy)
```bash
curl -X POST http://localhost:8080/trading/trade \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "INFY",
    "quantity": 75,
    "price": 62.00,
    "type": "BUY"
  }'
```

### 7. Execute Trade (Sell)
```bash
curl -X POST http://localhost:8080/trading/trade \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "RELIANCE",
    "quantity": 25,
    "price": 102.50,
    "type": "SELL"
  }'
```

### 8. Deposit Cash
```bash
curl -X POST http://localhost:8080/trading/deposit \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000
  }'
```

### 9. Create Price Alert
```bash
curl -X POST http://localhost:8080/trading/alerts \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "RELIANCE",
    "price": 105.00,
    "condition": ">="
  }'
```

**Conditions:** `>`, `<`, `>=`, `<=`, `==`

### 10. Get All Alerts
```bash
curl -X GET http://localhost:8080/trading/alerts \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 11. Get User Profile
```bash
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Response:**
```json
{
  "id": "user123",
  "username": "harsh",
  "email": "harsh@example.com",
  "phone": "+91 98765 43210",
  "bank_account": "HDFC Bank - Savings | ***1234",
  "bank_account_status": "Verified",
  "member_since": "2025-03-30T00:00:00Z",
  "is_verified": true,
  "theme": "Light",
  "notification_alerts": true,
  "notification_trades": true,
  "notification_news": false,
  "two_factor_enabled": true
}
```

### 12. Update Profile Settings
```bash
curl -X PUT http://localhost:8080/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "Dark",
    "notification_alerts": true,
    "notification_trades": true,
    "notification_news": false
  }'
```

### 13. Get Active Sessions
```bash
curl -X GET http://localhost:8080/user/sessions \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 14. Logout
```bash
curl -X POST http://localhost:8080/user/logout \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## WebSocket (Real-time Price Stream)

### Connect to WebSocket:
```bash
# Using websocat (install: cargo install websocat)
websocat "ws://localhost:8080/ws?token=YOUR_TOKEN_HERE"

# Or using Node.js
node -e "
const ws = new (require('ws'))('ws://localhost:8080/ws?token=YOUR_TOKEN_HERE');
ws.on('message', data => console.log(JSON.parse(data)));
"
```

**Receives prices every 500ms:**
```json
{
  "symbol": "RELIANCE",
  "currentPrice": 102.50,
  "percentageChange": 1.25,
  "timestamp": "2026-03-30T10:30:00Z"
}
```

---

## Complete Testing Flow

```bash
# 1. Start server
cd StockTrack && go run ./cmd/server

# 2. In another terminal, execute these commands in order:

# Register
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"pass123"}' \
  | jq -r '.token')

echo "Your Token: $TOKEN"

# Get wallet
curl -X GET http://localhost:8080/trading/wallet \
  -H "Authorization: Bearer $TOKEN" | jq

# Deposit $50k
curl -X POST http://localhost:8080/trading/deposit \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":50000}' | jq

# Buy RELIANCE
curl -X POST http://localhost:8080/trading/trade \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"RELIANCE","quantity":100,"price":102.50,"type":"BUY"}' | jq

# Check wallet again
curl -X GET http://localhost:8080/trading/wallet \
  -H "Authorization: Bearer $TOKEN" | jq

# Create alert
curl -X POST http://localhost:8080/trading/alerts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"RELIANCE","price":105.00,"condition":">="}' | jq

# View alerts
curl -X GET http://localhost:8080/trading/alerts \
  -H "Authorization: Bearer $TOKEN" | jq

# Get profile
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## API Endpoint Summary

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/health` | ❌ | Health check |
| POST | `/auth/register` | ❌ | Create account |
| POST | `/auth/login` | ❌ | Get JWT token |
| GET | `/trading/wallet` | ✅ | Get portfolio |
| POST | `/trading/trade` | ✅ | Buy/Sell stocks |
| POST | `/trading/deposit` | ✅ | Add funds |
| GET | `/trading/candles` | ✅ | Get OHLC data |
| POST | `/trading/alerts` | ✅ | Create alert |
| GET | `/trading/alerts` | ✅ | Get alerts |
| GET | `/user/profile` | ✅ | User profile |
| PUT | `/user/profile` | ✅ | Update settings |
| GET | `/user/sessions` | ✅ | Active sessions |
| POST | `/user/logout` | ✅ | End session |
| WS | `/ws` | ✅ | Real-time prices |

---

## Popular Stock Symbols

Use these in your requests:

- **RELIANCE** - Reliance Industries
- **INFY** - Infosys
- **TCS** - Tata Consultancy Services
- **HDFC** - HDFC Bank
- **ICICI** - ICICI Bank
- **SBIN** - State Bank of India
- **LT** - Larsen & Toubro
- **BAJAJ** - Bajaj Auto
- **WIPRO** - Wipro
- **MARUTI** - Maruti Suzuki

---

## Notes

- Replace `YOUR_TOKEN_HERE` with the actual token from login
- Port: **8080** (development)
- All timestamps in ISO 8601 format (UTC)
- Error codes: 400 (bad request), 401 (unauthorized), 500 (server error)
- Maximum candle limit: 1000

---

## Visual API Documentation

Open this file in your browser for interactive docs:
📄 **api-docs.html** - Beautiful HTML API documentation with examples

---

**Last Updated:** March 30, 2026
