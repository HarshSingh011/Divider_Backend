# WebSocket Real Quantity Display - With User Holdings

## 🔍 Understanding the Socket Data for RELIANCE-CE-2900

When you hold **236 shares** of RELIANCE-CE-2900, here's what the socket broadcasts:

### ✅ Correct WebSocket Response (RELIANCE-CE-2900 with your 236 holding)

```json
{
  "symbol": "RELIANCE-CE-2900",
  "currentPrice": 45.75,
  "percentageChange": 0.549,
  "totalQuantity": 100000,
  "availableQuantity": 99764,
  "heldQuantity": 236,
  "timestamp": "2026-04-05T14:30:45.500Z"
}
```

---

## 📊 What Each Field Means

| Field | Value | Meaning |
|-------|-------|---------|
| **symbol** | RELIANCE-CE-2900 | Stock ticker |
| **currentPrice** | 45.75 | Latest market price |
| **percentageChange** | 0.549 | Price change % |
| **totalQuantity** | 100,000 | **Total stock available in market** |
| **heldQuantity** | 236 | **You own this many shares** |
| **availableQuantity** | 99,764 | **Market has this many left** (100,000 - 236) |

---

## 🎯 For Your UI Display

### When User Opens Dashboard

**Show on Home Screen:**
```
💰 Your Holdings in RELIANCE-CE-2900
├─ Quantity Held: 236 shares
│  (from heldQuantity field)
│
├─ Current Price: ₹45.75
│  (from currentPrice field)
│
├─ Market Availability: 99,764 shares available
│  (from availableQuantity field)
│
└─ Total Market: 100,000 shares exists
   (from totalQuantity field)
```

### When User Tries to Buy More

**Show this information:**
```
Max you can buy: 99,764 shares
(system will prevent buying above availableQuantity)

Example:
├─ Total Available in Market: 99,764
├─ You want to buy: 99,764? ✅ Allowed!
├─ You want to buy: 99,765? ❌ Not enough stock!
└─ Current Market Holds: 236 (yours)
```

---

## 📱 Complete Room Data Every 500ms

When socket broadcasts to your app, each company shows:

```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 45.75,
    "percentageChange": 0.549,
    "totalQuantity": 100000,
    "availableQuantity": 99764,
    "heldQuantity": 236,
    "timestamp": "2026-04-05T14:30:45.500Z"
  },
  {
    "symbol": "HDFC-PE-1400",
    "currentPrice": 12.35,
    "percentageChange": 1.230,
    "totalQuantity": 50000,
    "availableQuantity": 48500,
    "heldQuantity": 1500,
    "timestamp": "2026-04-05T14:30:45.500Z"
  },
  {
    "symbol": "INFY-CE-1500",
    "currentPrice": 62.50,
    "percentageChange": -0.398,
    "totalQuantity": 75000,
    "availableQuantity": 73200,
    "heldQuantity": 1800,
    "timestamp": "2026-04-05T14:30:45.500Z"
  },
  {
    "symbol": "TCS-PE-3500",
    "currentPrice": 85.60,
    "percentageChange": 0.352,
    "totalQuantity": 60000,
    "availableQuantity": 59900,
    "heldQuantity": 100,
    "timestamp": "2026-04-05T14:30:45.500Z"
  },
  {
    "symbol": "TCS-CE-4500",
    "currentPrice": 55.45,
    "percentageChange": 0.453,
    "totalQuantity": 80000,
    "availableQuantity": 80000,
    "heldQuantity": 0,
    "timestamp": "2026-04-05T14:30:45.500Z"
  },
  {
    "symbol": "ITC-PE-2800",
    "currentPrice": 38.90,
    "percentageChange": 0.387,
    "totalQuantity": 90000,
    "availableQuantity": 89750,
    "heldQuantity": 250,
    "timestamp": "2026-04-05T14:30:45.500Z"
  }
]
```

---

## 🛠️ JavaScript - How to Display This

```javascript
// When socket message arrives
socket.onmessage = function(event) {
  const marketData = JSON.parse(event.data);

  marketData.forEach(stock => {
    // Show available quantity (not total!)
    const displayText = `
      ${stock.symbol}
      ├─ You Own: ${stock.heldQuantity.toLocaleString()} shares
      ├─ Available to Buy: ${stock.availableQuantity.toLocaleString()} shares
      ├─ Price: ₹${stock.currentPrice.toFixed(2)}
      └─ Change: ${stock.percentageChange > 0 ? '+' : ''}${stock.percentageChange.toFixed(2)}%
    `;
    
    console.log(displayText);
    
    // Update UI
    document.getElementById(`stock-${stock.symbol}-available`).textContent = 
      stock.availableQuantity.toLocaleString();
    
    document.getElementById(`stock-${stock.symbol}-held`).textContent = 
      stock.heldQuantity.toLocaleString();
  });
};
```

---

## ✅ Market Board Example

### Your Market Board Should Show:

```
┌─────────────────────────────────────────────────────────────┐
│                    Market Board                             │
├──────────────────────────────────────────────────────────────┤
│ Symbol             │ Price  │ Available │ You Own │ Change  │
├──────────────────────────────────────────────────────────────┤
│ RELIANCE-CE-2900   │ ₹45.75 │ 99,764    │ 236     │ +0.55%  │
│ HDFC-PE-1400       │ ₹12.35 │ 48,500    │ -       │ +1.23%  │
│ INFY-CE-1500       │ ₹62.50 │ 73,200    │ -       │ -0.40%  │
│ TCS-PE-3500        │ ₹85.60 │ 59,900    │ -       │ +0.35%  │
│ TCS-CE-4500        │ ₹55.45 │ 80,000    │ -       │ +0.45%  │
│ ITC-PE-2800        │ ₹38.90 │ 89,750    │ -       │ +0.39%  │
└──────────────────────────────────────────────────────────────┘
```

---

## 📋 HTML Example for Market Board

```html
<table class="market-board">
  <thead>
    <tr>
      <th>Symbol</th>
      <th>Price</th>
      <th>Available in Market</th>
      <th>You Own</th>
      <th>% Change</th>
      <th>Action</th>
    </tr>
  </thead>
  <tbody id="marketBody">
    <!-- Populated by JavaScript -->
  </tbody>
</table>

<script>
socket.onmessage = function(event) {
  const marketData = JSON.parse(event.data);
  const tbody = document.getElementById('marketBody');
  tbody.innerHTML = '';

  marketData.forEach(stock => {
    const row = `
      <tr>
        <td><strong>${stock.symbol}</strong></td>
        <td>₹${stock.currentPrice.toFixed(2)}</td>
        <td>
          <span class="${stock.availableQuantity < 10000 ? 'warning' : 'normal'}">
            ${stock.availableQuantity.toLocaleString()}
          </span>
        </td>
        <td>${stock.heldQuantity > 0 ? stock.heldQuantity.toLocaleString() : '-'}</td>
        <td class="${stock.percentageChange > 0 ? 'green' : 'red'}">
          ${stock.percentageChange > 0 ? '+' : ''}${stock.percentageChange.toFixed(2)}%
        </td>
        <td>
          <button onclick="buyStock('${stock.symbol}', ${stock.availableQuantity})">
            Buy
          </button>
        </td>
      </tr>
    `;
    tbody.innerHTML += row;
  });
};
</script>
```

---

## 🔐 Buy Button Logic

When user clicks BUY, use `availableQuantity` to validate:

```javascript
function buyStock(symbol, availableQty) {
  const userQuantity = prompt(`How many shares? (Max: ${availableQty.toLocaleString()})`);
  
  if (userQuantity > availableQty) {
    alert(`❌ Only ${availableQty.toLocaleString()} available!`);
    return;
  }

  if (userQuantity <= 0) {
    alert('❌ Enter valid quantity');
    return;
  }

  // Proceed with buy
  executeBuy(symbol, userQuantity);
}
```

---

## 📊 Real-Time Scenario

### Initial Market State
```
RELIANCE-CE-2900
├─ Total in market: 100,000 shares
├─ Available: 100,000 shares
└─ Held: 0 shares (market just opened)
```

### After User (you) Buys 236 Shares
```
RELIANCE-CE-2900
├─ Total in market: 100,000 shares (unchanged)
├─ Available: 99,764 shares (100,000 - 236) ✅ REDUCED
└─ Held: 236 shares (your + all others)
```

### When Another User Tries to Buy 99,765
```
Request: Buy 99,765 shares
Available: 99,764 shares

❌ ERROR: Cannot buy! Only 99,764 available (you're holding 236)
```

---

## 🎯 Key Points

✅ **Socket broadcasts REAL available quantity** (total - all holdings)
✅ **Updates every 500ms** as traders buy/sell
✅ **Your 236 shares ARE being deducted** from the market
✅ **Display availableQuantity** on your market board
✅ **Use availableQuantity** for buy validation
✅ **Show heldQuantity** so users see their holdings

---

## 💡 Formula Behind the Scenes

```
availableQuantity = totalQuantity - heldQuantity

Where:
- totalQuantity = 100,000 (fixed per stock)
- heldQuantity = Sum of all users' net holdings
  = (User1 bought 100 - sold 10)
  + (User2 bought 200 - sold 50)
  + (You bought 236 - sold 0)
  + ...
  = Total held by all traders

Example with you holding 236:
- All others holding: 0 (in this example)
- You holding: 236
- heldQuantity = 236
- availableQuantity = 100,000 - 236 = 99,764
```

---

## ✅ Verification Checklist

- [ ] Socket shows availableQuantity (not totalQuantity) for trading
- [ ] availableQuantity = totalQuantity - heldQuantity
- [ ] When you buy 236, availableQuantity becomes 99,764 ✅
- [ ] Market board displays availableQuantity to users
- [ ] Buy validation checks against availableQuantity
- [ ] Updates happen every 500ms in real-time
- [ ] Other users see your 236 deducted from available

🚀 **Your 236 shares ARE being subtracted properly!**
