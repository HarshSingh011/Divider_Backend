# WebSocket JSON Message Examples

## Complete Real-Time Market Data

### Full WebSocket Broadcast (Every 500ms)

The server sends this JSON array to all connected clients every 500 milliseconds:

```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 46.23,
    "percentageChange": 1.602,
    "totalQuantity": 100000,
    "availableQuantity": 95000,
    "heldQuantity": 5000,
    "timestamp": "2026-04-05T14:30:45.123Z"
  },
  {
    "symbol": "HDFC-PE-1400",
    "currentPrice": 12.45,
    "percentageChange": 2.049,
    "totalQuantity": 50000,
    "availableQuantity": 48500,
    "heldQuantity": 1500,
    "timestamp": "2026-04-05T14:30:45.123Z"
  },
  {
    "symbol": "INFY-CE-1500",
    "currentPrice": 63.10,
    "percentageChange": 0.556,
    "totalQuantity": 75000,
    "availableQuantity": 73200,
    "heldQuantity": 1800,
    "timestamp": "2026-04-05T14:30:45.123Z"
  },
  {
    "symbol": "TCS-PE-3500",
    "currentPrice": 85.95,
    "percentageChange": 0.761,
    "totalQuantity": 60000,
    "availableQuantity": 58900,
    "heldQuantity": 1100,
    "timestamp": "2026-04-05T14:30:45.123Z"
  },
  {
    "symbol": "TCS-CE-4500",
    "currentPrice": 55.87,
    "percentageChange": 1.218,
    "totalQuantity": 80000,
    "availableQuantity": 80000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T14:30:45.123Z"
  },
  {
    "symbol": "ITC-PE-2800",
    "currentPrice": 39.12,
    "percentageChange": 0.953,
    "totalQuantity": 90000,
    "availableQuantity": 89750,
    "heldQuantity": 250,
    "timestamp": "2026-04-05T14:30:45.123Z"
  }
]
```

---

## Real-Time Updates (500ms Intervals)

### Update 1 (Time: 10:30:45.000Z)

```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 45.50,
    "percentageChange": 0.000,
    "totalQuantity": 100000,
    "availableQuantity": 100000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  },
  {
    "symbol": "HDFC-PE-1400",
    "currentPrice": 12.20,
    "percentageChange": 0.000,
    "totalQuantity": 50000,
    "availableQuantity": 50000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  },
  {
    "symbol": "INFY-CE-1500",
    "currentPrice": 62.75,
    "percentageChange": 0.000,
    "totalQuantity": 75000,
    "availableQuantity": 75000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  },
  {
    "symbol": "TCS-PE-3500",
    "currentPrice": 85.30,
    "percentageChange": 0.000,
    "totalQuantity": 60000,
    "availableQuantity": 60000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  },
  {
    "symbol": "TCS-CE-4500",
    "currentPrice": 55.20,
    "percentageChange": 0.000,
    "totalQuantity": 80000,
    "availableQuantity": 80000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  },
  {
    "symbol": "ITC-PE-2800",
    "currentPrice": 38.75,
    "percentageChange": 0.000,
    "totalQuantity": 90000,
    "availableQuantity": 90000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.000Z"
  }
]
```

---

### Update 2 (Time: 10:30:45.500Z) - After First Transaction

**Scenario**: User A bought 5,000 shares of RELIANCE-CE-2900

```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 45.75,
    "percentageChange": 0.549,
    "totalQuantity": 100000,
    "availableQuantity": 95000,
    "heldQuantity": 5000,
    "timestamp": "2026-04-05T10:30:45.500Z"
  },
  {
    "symbol": "HDFC-PE-1400",
    "currentPrice": 12.35,
    "percentageChange": 1.230,
    "totalQuantity": 50000,
    "availableQuantity": 50000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.500Z"
  },
  {
    "symbol": "INFY-CE-1500",
    "currentPrice": 62.50,
    "percentageChange": -0.398,
    "totalQuantity": 75000,
    "availableQuantity": 75000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.500Z"
  },
  {
    "symbol": "TCS-PE-3500",
    "currentPrice": 85.60,
    "percentageChange": 0.352,
    "totalQuantity": 60000,
    "availableQuantity": 60000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.500Z"
  },
  {
    "symbol": "TCS-CE-4500",
    "currentPrice": 55.45,
    "percentageChange": 0.453,
    "totalQuantity": 80000,
    "availableQuantity": 80000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.500Z"
  },
  {
    "symbol": "ITC-PE-2800",
    "currentPrice": 38.90,
    "percentageChange": 0.387,
    "totalQuantity": 90000,
    "availableQuantity": 90000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:45.500Z"
  }
]
```

---

### Update 3 (Time: 10:30:46.000Z) - User Sells

**Scenario**: User A sold 3,000 of their 5,000 RELIANCE shares

```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 46.00,
    "percentageChange": 1.099,
    "totalQuantity": 100000,
    "availableQuantity": 98000,
    "heldQuantity": 2000,
    "timestamp": "2026-04-05T10:30:46.000Z"
  },
  {
    "symbol": "HDFC-PE-1400",
    "currentPrice": 12.50,
    "percentageChange": 2.459,
    "totalQuantity": 50000,
    "availableQuantity": 49500,
    "heldQuantity": 500,
    "timestamp": "2026-04-05T10:30:46.000Z"
  },
  {
    "symbol": "INFY-CE-1500",
    "currentPrice": 62.80,
    "percentageChange": -0.079,
    "totalQuantity": 75000,
    "availableQuantity": 73000,
    "heldQuantity": 2000,
    "timestamp": "2026-04-05T10:30:46.000Z"
  },
  {
    "symbol": "TCS-PE-3500",
    "currentPrice": 85.80,
    "percentageChange": 0.585,
    "totalQuantity": 60000,
    "availableQuantity": 59500,
    "heldQuantity": 500,
    "timestamp": "2026-04-05T10:30:46.000Z"
  },
  {
    "symbol": "TCS-CE-4500",
    "currentPrice": 55.70,
    "percentageChange": 0.905,
    "totalQuantity": 80000,
    "availableQuantity": 80000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T10:30:46.000Z"
  },
  {
    "symbol": "ITC-PE-2800",
    "currentPrice": 39.20,
    "percentageChange": 1.165,
    "totalQuantity": 90000,
    "availableQuantity": 89500,
    "heldQuantity": 500,
    "timestamp": "2026-04-05T10:30:46.000Z"
  }
]
```

---

## Individual Stock Examples

### RELIANCE-CE-2900 (Starting State)

```json
{
  "symbol": "RELIANCE-CE-2900",
  "currentPrice": 45.50,
  "percentageChange": 0.0,
  "totalQuantity": 100000,
  "availableQuantity": 100000,
  "heldQuantity": 0,
  "timestamp": "2026-04-05T10:30:45.000Z"
}
```

### RELIANCE-CE-2900 (After 50% Market Saturation)

```json
{
  "symbol": "RELIANCE-CE-2900",
  "currentPrice": 47.25,
  "percentageChange": 3.846,
  "totalQuantity": 100000,
  "availableQuantity": 50000,
  "heldQuantity": 50000,
  "timestamp": "2026-04-05T10:35:15.000Z"
}
```

### TCS-CE-4500 (Newly Listed - Fresh)

```json
{
  "symbol": "TCS-CE-4500",
  "currentPrice": 55.20,
  "percentageChange": 0.0,
  "totalQuantity": 80000,
  "availableQuantity": 80000,
  "heldQuantity": 0,
  "timestamp": "2026-04-05T10:30:45.000Z"
}
```

### TCS-CE-4500 (After Some Trading)

```json
{
  "symbol": "TCS-CE-4500",
  "currentPrice": 56.10,
  "percentageChange": 1.630,
  "totalQuantity": 80000,
  "availableQuantity": 75500,
  "heldQuantity": 4500,
  "timestamp": "2026-04-05T10:45:20.000Z"
}
```

### ITC-PE-2800 (Low Availability Alert)

```json
{
  "symbol": "ITC-PE-2800",
  "currentPrice": 38.95,
  "percentageChange": 0.649,
  "totalQuantity": 90000,
  "availableQuantity": 2500,
  "heldQuantity": 87500,
  "timestamp": "2026-04-05T11:00:00.000Z"
}
```

---

## JavaScript: How to Parse WebSocket Messages

```javascript
// Connect to WebSocket
const socket = new WebSocket('ws://localhost:8080/ws');

// Set headers for authentication
const userHeaders = {
  'X-User-ID': 'user_12345',
  'X-Email': 'john@example.com',
  'X-Username': 'john_doe'
};

// Handle incoming messages
socket.onmessage = function(event) {
  // Parse the JSON array
  const marketData = JSON.parse(event.data);
  
  console.log('📊 MARKET UPDATE - ' + new Date().toLocaleTimeString());
  console.log('================================================');
  
  // Process each stock
  marketData.forEach(stock => {
    const priceColor = stock.percentageChange > 0 ? '📈' : '📉';
    const availPercent = (stock.availableQuantity / stock.totalQuantity * 100).toFixed(1);
    const availStatus = availPercent > 80 ? '✅ HIGH' : availPercent > 50 ? '⚠️  MEDIUM' : '🔴 LOW';
    
    console.log(`
    ${stock.symbol}
    ├─ Price: ₹${stock.currentPrice.toFixed(2)} ${priceColor} ${stock.percentageChange.toFixed(2)}%
    ├─ Available: ${stock.availableQuantity.toLocaleString()} / ${stock.totalQuantity.toLocaleString()}
    ├─ Status: ${availStatus}
    └─ Update: ${new Date(stock.timestamp).toLocaleTimeString()}
    `);
  });
};

// Display real-time prices in a table
socket.onmessage = function(event) {
  const marketData = JSON.parse(event.data);
  
  console.table(marketData.map(stock => ({
    Symbol: stock.symbol,
    Price: '₹' + stock.currentPrice.toFixed(2),
    'Change %': stock.percentageChange.toFixed(2) + '%',
    'Available': stock.availableQuantity.toLocaleString(),
    'Total': stock.totalQuantity.toLocaleString(),
    'Held': stock.heldQuantity.toLocaleString(),
    'Time': new Date(stock.timestamp).toLocaleTimeString()
  })));
};

// Error handling
socket.onerror = function(error) {
  console.error('WebSocket Error:', error);
};

socket.onclose = function() {
  console.log('WebSocket connection closed');
};
```

---

## React Component Example

```jsx
import React, { useEffect, useState } from 'react';

function MarketBoard() {
  const [marketData, setMarketData] = useState([]);
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    // Connect to WebSocket
    const ws = new WebSocket('ws://localhost:8080/ws');
    
    ws.onopen = () => {
      console.log('WebSocket connected');
      setSocket(ws);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setMarketData(data);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => ws.close();
  }, []);

  return (
    <div className="market-board">
      <h1>📊 Market Data - Real-Time Updates</h1>
      <table>
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Price</th>
            <th>Change %</th>
            <th>Available</th>
            <th>Total</th>
            <th>Held</th>
            <th>Availability</th>
          </tr>
        </thead>
        <tbody>
          {marketData.map((stock) => {
            const availPercent = (stock.availableQuantity / stock.totalQuantity * 100).toFixed(1);
            const isPositive = stock.percentageChange > 0;
            
            return (
              <tr key={stock.symbol} className={isPositive ? 'positive' : 'negative'}>
                <td><strong>{stock.symbol}</strong></td>
                <td>₹{stock.currentPrice.toFixed(2)}</td>
                <td className={isPositive ? 'green' : 'red'}>
                  {isPositive ? '📈' : '📉'} {stock.percentageChange.toFixed(2)}%
                </td>
                <td>{stock.availableQuantity.toLocaleString()}</td>
                <td>{stock.totalQuantity.toLocaleString()}</td>
                <td>{stock.heldQuantity.toLocaleString()}</td>
                <td>
                  <div className="progress-bar">
                    <div 
                      className="progress" 
                      style={{width: availPercent + '%'}}
                    >
                      {availPercent}%
                    </div>
                  </div>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default MarketBoard;
```

---

## Python Example

```python
import websocket
import json
import threading
from datetime import datetime

class MarketDataClient:
    def __init__(self, url='ws://localhost:8080/ws'):
        self.url = url
        self.ws = None
        self.market_data = []
    
    def on_message(self, ws, message):
        """Handle incoming WebSocket messages"""
        self.market_data = json.loads(message)
        self.display_market_data()
    
    def on_error(self, ws, error):
        """Handle WebSocket errors"""
        print(f"Error: {error}")
    
    def on_close(self, ws):
        """Handle WebSocket close"""
        print("WebSocket connection closed")
    
    def on_open(self, ws):
        """Handle WebSocket open"""
        print("WebSocket connected!")
    
    def display_market_data(self):
        """Display market data in console"""
        print("\n📊 MARKET UPDATE - " + datetime.now().strftime("%H:%M:%S"))
        print("=" * 80)
        
        for stock in self.market_data:
            available_percent = (stock['availableQuantity'] / stock['totalQuantity'] * 100)
            price_indicator = "📈" if stock['percentageChange'] > 0 else "📉"
            status = "✅ HIGH" if available_percent > 80 else "⚠️ MEDIUM" if available_percent > 50 else "🔴 LOW"
            
            print(f"""
{stock['symbol']}
├─ Price: ₹{stock['currentPrice']:.2f} {price_indicator} {stock['percentageChange']:.2f}%
├─ Available: {stock['availableQuantity']:,} / {stock['totalQuantity']:,}
├─ Held: {stock['heldQuantity']:,}
├─ Availability: {status} ({available_percent:.1f}%)
└─ Time: {stock['timestamp']}
""")
    
    def connect(self):
        """Connect to WebSocket"""
        websocket.enableTrace(False)
        self.ws = websocket.WebSocketApp(
            self.url,
            on_message=self.on_message,
            on_error=self.on_error,
            on_close=self.on_close
        )
        self.ws.on_open = self.on_open
        
        # Run in background thread
        wst = threading.Thread(target=self.ws.run_forever)
        wst.daemon = True
        wst.start()

# Usage
if __name__ == "__main__":
    client = MarketDataClient()
    client.connect()
    
    # Keep running
    try:
        while True:
            pass
    except KeyboardInterrupt:
        print("Closing connection...")
        if client.ws:
            client.ws.close()
```

---

## Message Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    WebSocket Server                             │
│                 Broadcasting Market Data                        │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                    Every 500ms │ Broadcast
                               │
                 ┌─────────────┼─────────────┐
                 │             │             │
                 ▼             ▼             ▼
         ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
         │   Client 1   │ │   Client 2   │ │   Client N   │
         │  (User A)    │ │  (User B)    │ │  (User Z)    │
         └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
                │                │                │
                │ Parse JSON     │ Parse JSON     │ Parse JSON
                │                │                │
                ▼                ▼                ▼
         ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
         │   Display    │ │   Display    │ │   Display    │
         │  Market UI   │ │  Market UI   │ │  Market UI   │
         └──────────────┘ └──────────────┘ └──────────────┘
```

---

## Summary

| Aspect | Details |
|--------|---------|
| **Message Type** | JSON Array |
| **Elements** | 6 stocks (RELIANCE, HDFC, INFY, TCS-PE, TCS-CE, ITC) |
| **Frequency** | Every 500ms |
| **Fields** | Symbol, Price, % Change, Total Qty, Available Qty, Held Qty, Timestamp |
| **Updates** | Real-time as users trade |
| **Connection** | WebSocket (ws://) |
| **Bandwidth** | ~800 bytes/sec per client |

Each message contains **complete market state** for all 6 companies! 🚀
