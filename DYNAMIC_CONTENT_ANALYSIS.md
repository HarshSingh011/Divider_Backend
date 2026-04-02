# App.html - Dynamic Content Analysis

## ✅ FULLY DYNAMIC (Loaded from API)

### 1. **Home Screen** ✅
- **Balance Display**: Loads from `/trading/wallet` API
  - `home-balance`: Shows total cash balance (₹)
  - `home-change`: Shows portfolio change percentage
- **Live Prices**: Real-time WebSocket connection
  - Connects to `/ws` endpoint
  - Updates `prices-container` with live market data
- **Wallet Data**: Fetches user's cash and holdings from API
- **User Profile**: Loads basic profile from `/user/profile` API

### 2. **Charts Screen** ✅
- **Candle Data**: Loads from `/trading/candles` API
  - Fetches OHLC data for selected symbol
  - Displays in chart container
  - Updates dynamically when symbol changes

### 3. **Portfolio Screen** ✅
- **Holdings**: Loads from `/trading/wallet` API
  - Shows all user holdings (symbol, quantity, value, P&L)
  - Updates dynamically
- **Transactions**: Loads from `/trading/wallet` API
  - Shows all buy/sell transactions
  - Marks as BUY 🟢 or SELL 🔴

### 4. **Alerts Screen** ✅
- **Active Alerts**: Loads from `/trading/alerts` API (GET)
  - Lists all active price alerts
  - Shows condition (>, <, ≥, ≤) and threshold price
- **Create Alert**: POST to `/trading/alerts` API
  - Creates new alert with symbol, condition, price
  - Refreshes alert list after creation

### 5. **Profile Screen** ✅ (FIXED)
- **User Profile**: Loads from `/user/profile` API
  - Username (or email if no username)
  - Email address
  - Portfolio value
  - Member since date
  - Account status
  - Logout button functionality

## 📊 DATA LOADING FUNCTIONS

| Function | Endpoint | Method | Dynamic |
|----------|----------|--------|---------|
| `loadWalletData()` | `/trading/wallet` | GET | ✅ Yes |
| `loadCandleData(symbol)` | `/trading/candles` | GET | ✅ Yes |
| `loadAlertsData()` | `/trading/alerts` | GET | ✅ Yes |
| `createAlert()` | `/trading/alerts` | POST | ✅ Yes |
| `loadUserProfile()` | `/user/profile` | GET | ✅ Yes |
| `connectWebSocket()` | `/ws` | WS | ✅ Yes (Real-time) |

## 🔌 API ENDPOINTS USED

### Authentication
- `/auth/login` - POST (Login)
- `/auth/register` - POST (Register)

### Trading
- `/trading/wallet` - GET (Fetch wallet & holdings)
- `/trading/deposit` - POST (Deposit cash)
- `/trading/trade` - POST (Execute buy/sell)
- `/trading/candles` - GET (Fetch OHLC data)
- `/trading/alerts` - GET/POST (Get/create alerts)

### User Profile
- `/user/profile` - GET (Fetch user profile)
- `/user/sessions` - GET (User sessions)
- `/user/logout` - POST (Logout)

### WebSocket
- `/ws` - WS (Real-time market data)

### Health
- `/health` - GET (API health check)

## 📱 SCREEN UPDATE MECHANISMS

### Trigger 1: Navigation
```javascript
switchScreen(index) → Triggers loadUserProfile() when profile screen opened
```

### Trigger 2: WebSocket
```javascript
connectWebSocket() → Real-time price updates via message handler
```

### Trigger 3: Manual Actions
```javascript
- Buying/Selling triggers loadWalletData()
- Creating Alert triggers loadAlertsData()
- Depositing Cash triggers loadWalletData()
```

## 🎯 ALL UI UPDATES ARE DYNAMIC

### Container IDs That Update from API Data:
- `home-balance` ← Wallet balance
- `home-change` ← Portfolio change
- `prices-container` ← WebSocket prices
- `portfolio-header-info` ← Holdings count
- `holdings-container` ← Holding details
- `transactions-container` ← Transaction history
- `alerts-container` ← Alert listings
- `profile-container` ← User profile info

## ⚠️ IMPORTANT NOTES

1. **No Hardcoded Static Data**: All content loads from API (except UI template structure)
2. **Error Handling**: Each API call has try-catch with user-friendly error messages
3. **Loading States**: Shows "Loading..." while fetching data
4. **Token-Based Authentication**: All protected endpoints require Bearer token
5. **Real-Time Updates**: WebSocket connected automatically after login
6. **Database**: All data persisted in Neon PostgreSQL (not in-memory)

## ✅ CONFIRMATION

**EVERYTHING IN APP.HTML IS FULLY DYNAMIC** - No static dummy data is displayed after login. All content is fetched from the backend API in real-time.

---
**Last Updated**: April 2, 2026
**API Backend**: https://divider-backend.onrender.com
**Database**: Neon PostgreSQL
