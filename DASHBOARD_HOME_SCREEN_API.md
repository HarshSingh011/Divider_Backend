# Dashboard Home Screen APIs

## 📊 Overview

For showing **Profit/Loss** and **Total Portfolio Value** on the home screen, you need to call ONE main API:

### **`GET /trading/wallet`** — Main API for Dashboard

This single endpoint returns everything you need for the home screen dashboard.

---

## 🎯 API Endpoint

```
GET http://localhost:8080/trading/wallet
```

### Headers Required

```
Authorization: Bearer <token>
X-User-ID: user_12345
X-Email: user@example.com
X-Username: john_doe
Content-Type: application/json
```

---

## 📤 Response Format

```json
{
  "success": true,
  "data": {
    "user_id": "user_12345",
    "total_balance": 125500.50,
    "available_cash": 45000.25,
    "invested_amount": 80500.25,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 100,
        "average_cost": 2500,
        "current_price": 2750,
        "unrealized_pnl": 25000,
        "percentage": 24.50
      },
      "TCS-PE-3500": {
        "symbol": "TCS-PE-3500",
        "quantity": 50,
        "average_cost": 3000,
        "current_price": 2950,
        "unrealized_pnl": -2500,
        "percentage": -8.33
      },
      "HDFC-PE-1400": {
        "symbol": "HDFC-PE-1400",
        "quantity": 200,
        "average_cost": 1400,
        "current_price": 1500,
        "unrealized_pnl": 20000,
        "percentage": 7.14
      }
    },
    "last_updated": "2026-04-05T14:30:45.123Z"
  }
}
```

---

## 💡 What Each Field Means (For Dashboard)

### **Top-Level Values**

| Field | Meaning | Display On |
|-------|---------|-----------|
| **total_balance** | Total portfolio value (cash + holdings at current prices) | 💰 **Total Portfolio Value** |
| **available_cash** | Cash ready to trade (not invested) | 💵 **Available Cash** |
| **invested_amount** | Money currently in stock positions | 📈 **Invested Amount** |
| **last_updated** | When this snapshot was calculated | ⏰ **Last Updated** |

### **Positions (Holdings)**

Each position shows individual stock performance:

| Field | Meaning | Display On |
|-------|---------|-----------|
| **symbol** | Stock ticker | Stock Name |
| **quantity** | How many shares you own | Shares Owned |
| **average_cost** | Average price you paid per share | Cost Basis |
| **current_price** | Latest market price | Current Price |
| **unrealized_pnl** | **Profit/Loss on THIS stock** | 📊 **P&L** (Color: Green if +, Red if -) |
| **percentage** | **Profit/Loss percentage** | 📈 **% Change** |

---

## 🎨 Dashboard Display Example

### HTML Layout

```html
<div class="dashboard">
  <!-- Top Cards -->
  <div class="card-row">
    <div class="card total-value">
      <h3>Total Portfolio Value</h3>
      <h2 id="totalValue">₹125,500.50</h2>
      <p>Including cash and holdings</p>
    </div>
    
    <div class="card">
      <h3>Available Cash</h3>
      <h2 id="availableCash">₹45,000.25</h2>
      <p>Ready to invest</p>
    </div>
    
    <div class="card">
      <h3>Invested Amount</h3>
      <h2 id="investedAmount">₹80,500.25</h2>
      <p>In stock holdings</p>
    </div>
  </div>

  <!-- Profit/Loss Summary -->
  <div class="card profit-loss">
    <h3>Total Unrealized P&L</h3>
    <h2 id="totalPnL" class="green">+₹42,500.00</h2>
    <p id="totalPnLPercent" class="green">+34.07%</p>
  </div>

  <!-- Holdings Table -->
  <table class="holdings">
    <thead>
      <tr>
        <th>Symbol</th>
        <th>Quantity</th>
        <th>Avg Cost</th>
        <th>Current Price</th>
        <th>P&L</th>
        <th>%</th>
      </tr>
    </thead>
    <tbody id="holdingsBody">
      <!-- Will be populated by JavaScript -->
    </tbody>
  </table>
</div>
```

---

## 🔧 JavaScript Implementation

### 1. Fetch Wallet Data

```javascript
async function loadDashboard() {
  const userID = localStorage.getItem('userID');
  const token = localStorage.getItem('token');

  try {
    const response = await fetch('http://localhost:8080/trading/wallet', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': userID,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    if (data.success) {
      displayDashboard(data.data);
    } else {
      console.error('API Error:', data.error);
    }
  } catch (error) {
    console.error('Failed to load dashboard:', error);
  }
}

// Call on page load
document.addEventListener('DOMContentLoaded', loadDashboard);
```

### 2. Display Dashboard Data

```javascript
function displayDashboard(walletSnapshot) {
  const {
    total_balance,
    available_cash,
    invested_amount,
    positions,
    last_updated
  } = walletSnapshot;

  // Update top cards
  document.getElementById('totalValue').textContent = 
    '₹' + total_balance.toLocaleString('en-IN', { maximumFractionDigits: 2 });
  
  document.getElementById('availableCash').textContent = 
    '₹' + available_cash.toLocaleString('en-IN', { maximumFractionDigits: 2 });
  
  document.getElementById('investedAmount').textContent = 
    '₹' + invested_amount.toLocaleString('en-IN', { maximumFractionDigits: 2 });

  // Calculate total P&L
  let totalPnL = 0;
  Object.values(positions).forEach(pos => {
    totalPnL += pos.unrealized_pnl;
  });

  const totalPnLPercent = (totalPnL / (invested_amount || 1)) * 100;
  
  const pnlElement = document.getElementById('totalPnL');
  const pnlPercentElement = document.getElementById('totalPnLPercent');
  
  pnlElement.textContent = (totalPnL >= 0 ? '+' : '') + 
    '₹' + totalPnL.toLocaleString('en-IN', { maximumFractionDigits: 2 });
  
  pnlPercentElement.textContent = (totalPnLPercent >= 0 ? '+' : '') + 
    totalPnLPercent.toFixed(2) + '%';
  
  // Color coding: Green for profit, Red for loss
  if (totalPnL >= 0) {
    pnlElement.classList.add('green');
    pnlElement.classList.remove('red');
  } else {
    pnlElement.classList.add('red');
    pnlElement.classList.remove('green');
  }

  // Display holdings table
  const tbody = document.getElementById('holdingsBody');
  tbody.innerHTML = '';

  Object.values(positions).forEach(position => {
    const row = document.createElement('tr');
    const pnlClass = position.unrealized_pnl >= 0 ? 'green' : 'red';
    
    row.innerHTML = `
      <td><strong>${position.symbol}</strong></td>
      <td>${position.quantity.toLocaleString('en-IN')}</td>
      <td>₹${position.average_cost.toLocaleString('en-IN', { maximumFractionDigits: 2 })}</td>
      <td>₹${position.current_price.toLocaleString('en-IN', { maximumFractionDigits: 2 })}</td>
      <td class="${pnlClass}">
        ${position.unrealized_pnl >= 0 ? '+' : ''}₹${position.unrealized_pnl.toLocaleString('en-IN', { maximumFractionDigits: 2 })}
      </td>
      <td class="${pnlClass}">
        ${position.percentage >= 0 ? '+' : ''}${position.percentage.toFixed(2)}%
      </td>
    `;
    
    tbody.appendChild(row);
  });

  // Update timestamp
  const lastUpdated = new Date(last_updated);
  document.getElementById('lastUpdated').textContent = 
    `Last updated: ${lastUpdated.toLocaleTimeString('en-IN')}`;
}
```

### 3. Auto-Refresh Dashboard

```javascript
// Refresh every 5 seconds
setInterval(loadDashboard, 5000);

// Or with WebSocket for real-time prices
function setupRealtimeUpdates() {
  const ws = new WebSocket('ws://localhost:8080/ws');

  ws.onmessage = function(event) {
    // When market data updates, refresh wallet snapshot
    const marketData = JSON.parse(event.data);
    console.log('Market updated, refreshing dashboard...', marketData);
    loadDashboard(); // Refresh to get latest P&L
  };

  ws.onerror = function(error) {
    console.error('WebSocket error:', error);
  };
}

// Call on page load
setupRealtimeUpdates();
```

---

## 💻 React Component Example

```jsx
import React, { useEffect, useState } from 'react';

function Dashboard() {
  const [walletData, setWalletData] = useState(null);
  const [totalPnL, setTotalPnL] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchWalletData();
    
    // Refresh every 5 seconds
    const interval = setInterval(fetchWalletData, 5000);
    return () => clearInterval(interval);
  }, []);

  const fetchWalletData = async () => {
    const userID = localStorage.getItem('userID');
    const token = localStorage.getItem('token');

    try {
      const response = await fetch('http://localhost:8080/trading/wallet', {
        headers: {
          'Authorization': `Bearer ${token}`,
          'X-User-ID': userID
        }
      });

      const data = await response.json();
      if (data.success) {
        setWalletData(data.data);
        
        // Calculate total P&L
        const pnl = Object.values(data.data.positions || {})
          .reduce((sum, pos) => sum + pos.unrealized_pnl, 0);
        setTotalPnL(pnl);
        
        setLoading(false);
      }
    } catch (error) {
      console.error('Failed to fetch wallet data:', error);
    }
  };

  if (loading || !walletData) return <div>Loading...</div>;

  const totalPnLPercent = (totalPnL / walletData.invested_amount) * 100;
  const isProfitable = totalPnL >= 0;

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>

      {/* Portfolio Value Cards */}
      <div className="card-grid">
        <Card
          title="Total Portfolio Value"
          value={`₹${walletData.total_balance.toLocaleString('en-IN', { maximumFractionDigits: 2 })}`}
          color="blue"
        />
        <Card
          title="Available Cash"
          value={`₹${walletData.available_cash.toLocaleString('en-IN', { maximumFractionDigits: 2 })}`}
          color="green"
        />
        <Card
          title="Invested Amount"
          value={`₹${walletData.invested_amount.toLocaleString('en-IN', { maximumFractionDigits: 2 })}`}
          color="orange"
        />
      </div>

      {/* Total P&L Card */}
      <Card
        title="Total Unrealized P&L"
        value={`${isProfitable ? '+' : ''}₹${totalPnL.toLocaleString('en-IN', { maximumFractionDigits: 2 })}`}
        subtitle={`${isProfitable ? '+' : ''}${totalPnLPercent.toFixed(2)}%`}
        color={isProfitable ? 'green' : 'red'}
      />

      {/* Holdings Table */}
      <HoldingsTable positions={walletData.positions} />

      {/* Last Updated */}
      <p className="timestamp">
        📅 Last updated: {new Date(walletData.last_updated).toLocaleString('en-IN')}
      </p>
    </div>
  );
}

function Card({ title, value, subtitle, color }) {
  return (
    <div className={`card card-${color}`}>
      <h3>{title}</h3>
      <h2>{value}</h2>
      {subtitle && <p>{subtitle}</p>}
    </div>
  );
}

function HoldingsTable({ positions }) {
  return (
    <table className="holdings-table">
      <thead>
        <tr>
          <th>Symbol</th>
          <th>Qty</th>
          <th>Avg Cost</th>
          <th>Current</th>
          <th>P&L</th>
          <th>%</th>
        </tr>
      </thead>
      <tbody>
        {Object.values(positions || {}).map((pos) => (
          <tr key={pos.symbol}>
            <td><strong>{pos.symbol}</strong></td>
            <td>{pos.quantity.toLocaleString('en-IN')}</td>
            <td>₹{pos.average_cost.toFixed(2)}</td>
            <td>₹{pos.current_price.toFixed(2)}</td>
            <td className={pos.unrealized_pnl >= 0 ? 'green' : 'red'}>
              {pos.unrealized_pnl >= 0 ? '+' : ''}₹{pos.unrealized_pnl.toLocaleString('en-IN', { maximumFractionDigits: 2 })}
            </td>
            <td className={pos.percentage >= 0 ? 'green' : 'red'}>
              {pos.percentage >= 0 ? '+' : ''}{pos.percentage.toFixed(2)}%
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default Dashboard;
```

---

## 📐 Calculation Formulas

### **Profit/Loss for Each Stock**

```
Unrealized P&L = (Current Price - Average Cost) × Quantity

Example:
- You bought RELIANCE at ₹2500 (100 shares) = ₹250,000
- Current price is ₹2750
- P&L = (2750 - 2500) × 100 = ₹25,000 profit
- P&L % = (25,000 / 250,000) × 100 = 10% gain
```

### **Total Portfolio Value**

```
Total Balance = Available Cash + Sum of (Current Price × Quantity) for all holdings

Example:
- Available Cash: ₹45,000
- RELIANCE: 100 × ₹2750 = ₹275,000
- TCS: 50 × ₹3000 = ₹150,000
- HDFC: 200 × ₹1500 = ₹300,000
- Total = 45,000 + 275,000 + 150,000 + 300,000 = ₹770,000
```

### **Total Unrealized P&L**

```
Total P&L = Sum of all individual position P&Ls

Example:
- RELIANCE P&L: +₹25,000
- TCS P&L: -₹2,500
- HDFC P&L: +₹20,000
- Total P&L = 25,000 - 2,500 + 20,000 = ₹42,500
```

---

## 🔄 Combining with WebSocket (Optional)

For **real-time P&L updates**, combine this API with WebSocket:

```javascript
const socket = new WebSocket('ws://localhost:8080/ws');

socket.onmessage = (event) => {
  const marketData = JSON.parse(event.data);
  
  // Update with new prices (optional: update UI immediately)
  // Then fetch fresh wallet snapshot for accurate P&L
  loadDashboard();
};
```

---

## 🛠️ CSS Styling Example

```css
.dashboard {
  padding: 20px;
  background: #f5f5f5;
  border-radius: 8px;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card h2 {
  font-size: 28px;
  margin: 10px 0;
}

.card-green h2 { color: #10b981; }
.card-red h2 { color: #ef4444; }
.card-blue h2 { color: #3b82f6; }
.card-orange h2 { color: #f59e0b; }

.green { color: #10b981; }
.red { color: #ef4444; }

.holdings-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  margin-top: 20px;
}

.holdings-table th,
.holdings-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.holdings-table th {
  background: #f9fafb;
  font-weight: 600;
}

.holdings-table tr:hover {
  background: #f9fafb;
}
```

---

## 📝 Summary

| Need | API Call | Response Field |
|------|----------|-----------------|
| **Total Portfolio Value** | `GET /trading/wallet` | `total_balance` |
| **Available Cash** | `GET /trading/wallet` | `available_cash` |
| **Invested Amount** | `GET /trading/wallet` | `invested_amount` |
| **Profit/Loss per Stock** | `GET /trading/wallet` | `positions[].unrealized_pnl` |
| **Profit/Loss %** | `GET /trading/wallet` | `positions[].percentage` |
| **All Holdings** | `GET /trading/wallet` | `positions` |
| **Real-time Prices** | `WebSocket /ws` | Array of market ticks |

---

## ✅ Implementation Checklist

- [ ] Call `GET /trading/wallet` endpoint
- [ ] Parse response and extract wallet data
- [ ] Display **Total Portfolio Value** (total_balance)
- [ ] Display **Available Cash** (available_cash)
- [ ] Display **Invested Amount** (invested_amount)
- [ ] Calculate and show **Total P&L** (sum of unrealized_pnl)
- [ ] Show individual stock P&L and percentages
- [ ] Color code P&L: Green (+), Red (-)
- [ ] Display holdings table with all positions
- [ ] Add auto-refresh every 5 seconds
- [ ] Optional: Connect WebSocket for real-time updates
- [ ] Format currency values with Indian locale

🚀 **Your home screen dashboard is ready!**
