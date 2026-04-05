# JSON Response Bodies for All 9 Test Cases (UI Ready)

## Overview
All API responses now return JSON format with `success` flag and structured error/data objects for easy UI parsing.

---

## Test Case 1: Sell Non-Owned Shares ❌

### Request
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 2700,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_HOLDINGS",
    "message": "insufficient holdings: you have 0.00 shares of RELIANCE-CE-2900 but trying to sell 10.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 10,
      "price": 2700,
      "type": "SELL"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Sell
Error Code: INSUFFICIENT_HOLDINGS
Message: You have 0.00 shares of RELIANCE-CE-2900 but trying to sell 10.00
```

---

## Test Case 2: Sell More Than Owned ❌

### Step 1: Buy 10 shares (Success)
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 8000,
    "fee": 40,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 10,
    "price": 800,
    "type": "BUY"
  }
}
```

### Step 2: Try to sell 15 (Failure)
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 15,
  "price": 900,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_HOLDINGS",
    "message": "insufficient holdings: you have 10.00 shares of RELIANCE-CE-2900 but trying to sell 15.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 15,
      "price": 900,
      "type": "SELL"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Sell
Error Code: INSUFFICIENT_HOLDINGS
Message: You have 10.00 shares but trying to sell 15.00
Action: You can only sell up to 10 shares
```

---

## Test Case 3: Sell Exact Amount Owned ✅

### Step 1: Buy 10 shares
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 8000,
    "fee": 40,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 10,
    "price": 800,
    "type": "BUY"
  }
}
```

### Step 2: Sell all 10 shares ✅
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 900,
  "type": "SELL"
}
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 9000,
    "fee": 45,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 10,
    "price": 900,
    "type": "SELL"
  }
}
```

### UI Display Suggestion
```
✅ Sell Successful
Symbol: RELIANCE-CE-2900
Quantity: 10 shares
Price: ₹900 per share
Total: ₹9000
Fee: ₹45
Net Proceeds: ₹8955
```

---

## Test Case 4: Sell Partial Amount ✅

### Step 1: Buy 100 shares
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 100,
  "price": 100,
  "type": "BUY"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 10000,
    "fee": 50,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 100,
    "price": 100,
    "type": "BUY"
  }
}
```

### Step 2: Sell 60 shares ✅
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 60,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 7200,
    "fee": 36,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 60,
    "price": 120,
    "type": "SELL"
  }
}
```

### Step 3: Check wallet to see remaining holdings
```json
GET /trading/wallet
X-User-ID: test_user
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 100000,
    "available_cash": 99914,
    "invested_amount": 4000,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 40,
        "average_cost": 100,
        "current_price": 120,
        "unrealized_pnl": 800,
        "percentage": 20
      }
    }
  }
}
```

### UI Display Suggestion
```
✅ Sell Successful
Sold: 60 shares @ ₹120
Proceeds: ₹7200
Fee: ₹36
Net Proceeds: ₹7164

Remaining Holdings:
RELIANCE-CE-2900: 40 shares
Average Cost: ₹100
Current Price: ₹120
Unrealized P&L: +₹800 (+20%)
```

---

## Test Case 5: Cannot Sell After Holdings Empty ❌

### Step 1: Buy 10 shares
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 100,
  "type": "BUY"
}
```

### Step 2: Sell all 10 shares ✅
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 120,
  "type": "SELL"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 1200,
    "fee": 6,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 10,
    "price": 120,
    "type": "SELL"
  }
}
```

### Step 3: Try to sell again ❌
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 5,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_HOLDINGS",
    "message": "insufficient holdings: you have 0.00 shares of RELIANCE-CE-2900 but trying to sell 5.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 5,
      "price": 120,
      "type": "SELL"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Sell
Error Code: INSUFFICIENT_HOLDINGS
Message: You have 0.00 shares of RELIANCE-CE-2900
Action: No holdings to sell. Buy shares first to sell later.
```

---

## Test Case 6: Buy Insufficient Cash ❌

### Request
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 1000,
  "price": 200,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_CASH",
    "message": "insufficient cash: you have ₹100000.00 but need ₹201000.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 1000,
      "price": 200,
      "type": "BUY"
    }
  }
}
```

### Calculation Shown in Response
```
Stock Cost: 1000 × 200 = ₹200,000
Fee (0.5%): ₹1,000
Total Needed: ₹201,000
Your Cash: ₹100,000
Shortfall: ₹101,000
```

### UI Display Suggestion
```
❌ Insufficient Funds
Error Code: INSUFFICIENT_CASH
Need: ₹201,000.00
Have: ₹100,000.00
Shortfall: ₹101,000.00

Suggestion: Deposit more cash or reduce quantity to 500 shares
Maximum you can buy: 499 shares @ ₹200
```

---

## Test Case 7: Multiple Buys Then Sell ✅

### Step 1-3: Buy 30, 20, and 50 shares
```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 30,
  "price": 100,
  "type": "BUY"
}
```
Response: 200 OK ✅

```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 20,
  "price": 110,
  "type": "BUY"
}
```
Response: 200 OK ✅

```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 90,
  "type": "BUY"
}
```
Response: 200 OK ✅

### Step 4: Try to sell 150 (have 100) ❌
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 150,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_HOLDINGS",
    "message": "insufficient holdings: you have 100.00 shares of RELIANCE-CE-2900 but trying to sell 150.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 150,
      "price": 120,
      "type": "SELL"
    }
  }
}
```

### Step 5: Sell 75 shares ✅
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 75,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 9000,
    "fee": 45,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 75,
    "price": 120,
    "type": "SELL"
  }
}
```

### Step 6: Check wallet
```json
GET /trading/wallet
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 100000,
    "available_cash": 99920,
    "invested_amount": 2500,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 25,
        "average_cost": 103.10,
        "current_price": 120,
        "unrealized_pnl": 423.75,
        "percentage": 4.11
      }
    }
  }
}
```

### UI Display Suggestion
```
Summary of Multiple Trades:
- BUY: 30 @ ₹100 ✓
- BUY: 20 @ ₹110 ✓
- BUY: 50 @ ₹90 ✓
Total Holdings: 100 shares
Average Cost: ₹103.10

Attempted to SELL 150 ❌ (Only have 100)
Then SELL: 75 @ ₹120 ✓

Remaining: 25 shares
Current Value: ₹3000
Unrealized P&L: +₹423.75 (+4.11%)
```

---

## Test Case 8: Multiple Symbols Independent ✅

### Step 1: Buy RELIANCE
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 100,
  "type": "BUY"
}
```
Response: 200 OK ✅

### Step 2: Buy TCS
```json
POST /trading/trade

{
  "symbol": "TCS-CE-3500",
  "quantity": 30,
  "price": 150,
  "type": "BUY"
}
```
Response: 200 OK ✅

### Step 3: Sell 40 RELIANCE ✅
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 40,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 4800,
    "fee": 24,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 40,
    "price": 120,
    "type": "SELL"
  }
}
```

### Step 4: Try to sell 50 TCS (have 30) ❌
```json
POST /trading/trade

{
  "symbol": "TCS-CE-3500",
  "quantity": 50,
  "price": 170,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_HOLDINGS",
    "message": "insufficient holdings: you have 30.00 shares of TCS-CE-3500 but trying to sell 50.00",
    "details": {
      "symbol": "TCS-CE-3500",
      "quantity": 50,
      "price": 170,
      "type": "SELL"
    }
  }
}
```

### Step 5: Check wallet
```json
GET /trading/wallet
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 100000,
    "available_cash": 99921,
    "invested_amount": 8500,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 10,
        "average_cost": 100,
        "current_price": 120,
        "unrealized_pnl": 200,
        "percentage": 20
      },
      "TCS-CE-3500": {
        "symbol": "TCS-CE-3500",
        "quantity": 30,
        "average_cost": 150,
        "current_price": 170,
        "unrealized_pnl": 600,
        "percentage": 4
      }
    }
  }
}
```

### UI Display Suggestion
```
Portfolio:
┌─────────────────────────────────────────────┐
│ RELIANCE-CE-2900                            │
│ Holdings: 10 shares | Avg: ₹100            │
│ P&L: +₹200 (+20%)                          │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ TCS-CE-3500                                 │
│ Holdings: 30 shares | Avg: ₹150            │
│ P&L: +₹600 (+4%)                           │
└─────────────────────────────────────────────┘

Cash: ₹99,921
Total Balance: ₹100,000
```

---

## Test Case 9: Invalid Quantity or Price ❌

### Request 1: Zero Quantity
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 0,
  "price": 100,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_QUANTITY_OR_PRICE",
    "message": "quantity and price must be positive",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 0,
      "price": 100,
      "type": "BUY"
    }
  }
}
```

---

### Request 2: Negative Quantity
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": -50,
  "price": 100,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_QUANTITY_OR_PRICE",
    "message": "quantity and price must be positive",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": -50,
      "price": 100,
      "type": "BUY"
    }
  }
}
```

---

### Request 3: Zero Price
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 0,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_QUANTITY_OR_PRICE",
    "message": "quantity and price must be positive",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 50,
      "price": 0,
      "type": "BUY"
    }
  }
}
```

---

### Request 4: Negative Price
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": -100,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_QUANTITY_OR_PRICE",
    "message": "quantity and price must be positive",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 50,
      "price": -100,
      "type": "BUY"
    }
  }
}
```

### UI Display Suggestion
```
❌ Invalid Input
Error Code: INVALID_QUANTITY_OR_PRICE
Message: Quantity and price must be positive numbers

Validation Rules:
✓ Quantity > 0
✓ Price > 0
✗ Your input: quantity=-50, price=-100
```

---

## Summary: Error Codes for UI Implementation

| Error Code | HTTP Status | Meaning | UI Action |
|-----------|------------|---------|-----------|
| `INSUFFICIENT_HOLDINGS` | 400 | Can't sell more than you own | Show max available shares |
| `INSUFFICIENT_CASH` | 400 | Can't buy without funds | Show cash needed vs available |
| `INVALID_QUANTITY_OR_PRICE` | 400 | Negative/zero values | Validate input fields |
| `INVALID_TRADE_TYPE` | 400 | Trade type not BUY/SELL | Show error |
| `INVALID_REQUEST_BODY` | 400 | JSON parsing error | Check request format |
| `METHOD_NOT_ALLOWED` | 405 | Wrong HTTP method | Check API docs |
| `UNAUTHORIZED` | 401 | Missing X-User-ID header | Login/authenticate |
| `WALLET_SNAPSHOT_ERROR` | 500 | Server error | Retry or contact support |
| `DEPOSIT_FAILED` | 400 | Deposit validation failed | Check amount > 0 |

---

## Frontend Integration Example

```javascript
// Parse response and show to user
fetch('/trading/trade', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-User-ID': 'test_user'
  },
  body: JSON.stringify({
    symbol: 'RELIANCE-CE-2900',
    quantity: 50,
    price: 2500,
    type: 'SELL'
  })
})
.then(res => res.json())
.then(data => {
  if (data.success) {
    // Success case
    showSuccess(`✅ Sold ${data.data.quantity} shares`);
    showDetails({
      total: data.data.total,
      fee: data.data.fee,
      proceeds: data.data.total - data.data.fee
    });
  } else {
    // Error case
    const errorCode = data.error.code;
    const errorMsg = data.error.message;
    const details = data.error.details;
    
    showError(errorCode, errorMsg);
    
    // Custom handling based on error code
    if (errorCode === 'INSUFFICIENT_HOLDINGS') {
      const available = details.quantity - details.quantity; // simplified
      showWarning(`You can only sell up to ${available} shares`);
    } else if (errorCode === 'INSUFFICIENT_CASH') {
      showWarning(`You need ₹${details.amount} more`);
    }
  }
});
```

---

## Frontend Integration Example - Stock Availability

```javascript
// Handle stock availability errors gracefully
fetch('/trading/trade', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-User-ID': 'test_user'
  },
  body: JSON.stringify({
    symbol: 'RELIANCE-CE-2900',
    quantity: 50,
    price: 2500,
    type: 'BUY'
  })
})
.then(res => res.json())
.then(data => {
  if (!data.success && data.error.code === 'INSUFFICIENT_STOCK') {
    // Extract available quantity from error message
    const match = data.error.message.match(/only ([\d.]+) shares/);
    const availableQty = match ? parseFloat(match[1]) : 0;
    
    showError('Not Enough Stock in Market', 
      `Only ${availableQty.toLocaleString()} shares available`);
    
    if (availableQty > 0) {
      showOption(
        `Buy ${availableQty.toLocaleString()} shares instead?`,
        () => submitBuyOrder(availableQty)
      );
    } else {
      showSuggestion('Wait for other traders to sell shares back into market');
      showAlternative('Browse similar stocks that have more availability');
    }
  }
});

// Helper: Show market depth to user
function displayMarketDepth(symbol, totalMarketCap, currentlyHeld, available) {
  const percentageHeld = (currentlyHeld / totalMarketCap * 100).toFixed(1);
  
  console.log(`
📊 Market Depth: ${symbol}
├─ Total in Market: ${totalMarketCap.toLocaleString()} shares
├─ Currently Held: ${currentlyHeld.toLocaleString()} (${percentageHeld}%)
└─ Available to Buy: ${available.toLocaleString()}

${available === 0 ? '⚠️  MARKET SATURATED - No shares available' 
  : available < (totalMarketCap * 0.1) ? '⚠️  LOW AVAILABILITY - Limited shares remaining'
  : '✅ Good availability'}
  `);
}

// Usage
displayMarketDepth('RELIANCE-CE-2900', 100, 95, 5);
// Output:
// 📊 Market Depth: RELIANCE-CE-2900
// ├─ Total in Market: 100 shares
// ├─ Currently Held: 95 (95%)
// └─ Available to Buy: 5 shares
// 
// ⚠️  LOW AVAILABILITY - Limited shares remaining
```

---

# MAXIMUM QUANTITY LIMIT (10,000 Shares Per Trade)

## Overview
To prevent market manipulation and ensure stable trading, there's a **maximum limit of 10,000 shares per trade** (both BUY and SELL). This limit applies per transaction.

---

## Test Case 10: Buy Above Maximum Quantity ❌

### Request
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 15000,
  "price": 100,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "MAXIMUM_QUANTITY_EXCEEDED",
    "message": "maximum quantity exceeded: you can only buy/sell up to 10000 shares per transaction, but trying to BUY 15000",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 15000,
      "price": 100,
      "type": "BUY"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Buy (Quantity Limit Exceeded)
Error Code: MAXIMUM_QUANTITY_EXCEEDED
Message: Maximum 10,000 shares per transaction
Requested: 15,000 shares
Maximum Allowed: 10,000 shares
Excess: 5,000 shares

Action: To buy all 15,000 shares, please split into 2 transactions:
- Transaction 1: Buy 10,000 shares @ ₹100
- Transaction 2: Buy 5,000 shares @ ₹100
```

---

## Test Case 11: Buy Exactly At Maximum Limit ✅

### Request
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10000,
  "price": 10,
  "type": "BUY"
}
```

**Note**: Requires sufficient cash. For 10,000 shares @ ₹10 = ₹100,000 + 0.5% fee = ₹100,500

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 100000,
    "fee": 500,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 10000,
    "price": 10,
    "type": "BUY"
  }
}
```

### Calculation Shown
```
Max Allowed: 10,000 shares
Requested: 10,000 shares ✓
Status: APPROVED

Stock Cost: 10,000 × ₹10 = ₹100,000
Fee (0.5%): ₹500
Total Cost: ₹100,500
```

### UI Display Suggestion
```
✅ Maximum Quantity Buy Successful
Symbol: RELIANCE-CE-2900
Quantity: 10,000 shares (At Limit)
Price: ₹10 per share
Total: ₹100,000
Fee: ₹500
Total Debit: ₹100,500

⚠️ Note: This is the maximum quantity allowed per transaction
```

---

## Test Case 12: Sell Above Maximum Quantity ❌

### Scenario
User tries to sell 15,000 shares (above limit)

### Request
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 15000,
  "price": 120,
  "type": "SELL"
}
```

### JSON Response (400 Bad Request)
```json
{
  "success": false,
  "error": {
    "code": "MAXIMUM_QUANTITY_EXCEEDED",
    "message": "maximum quantity exceeded: you can only buy/sell up to 10000 shares per transaction, but trying to SELL 15000",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 15000,
      "price": 120,
      "type": "SELL"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Sell (Quantity Limit Exceeded)
Error Code: MAXIMUM_QUANTITY_EXCEEDED
Message: Maximum 10,000 shares per transaction
Requested: 15,000 shares
Maximum Allowed: 10,000 shares

Action: Please split into 2 sell transactions:
- Transaction 1: Sell 10,000 shares @ ₹120
- Transaction 2: Sell 5,000 shares @ ₹120
```

---

## Test Case 13: Quantities Around The Limit

### Allowed Quantities ✅
- **100 shares** → ✅ Approved
- **5,000 shares** → ✅ Approved
- **9,999 shares** → ✅ Approved (just below limit)
- **10,000 shares** → ✅ Approved (exactly at limit)

### Rejected Quantities ❌
- **10,001 shares** → ❌ Rejected (1 over limit)
- **15,000 shares** → ❌ Rejected
- **50,000 shares** → ❌ Rejected (way over limit)

### Example: 9,999 Shares (Just Below Limit)
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 9999,
  "price": 100,
  "type": "BUY"
}
```

**Response (200 OK)** - Approved ✅
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 999900,
    "fee": 4999.50,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 9999,
    "price": 100,
    "type": "BUY"
  }
}
```

### Example: 10,001 Shares (Just Over Limit)
```json
POST /trading/trade

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10001,
  "price": 100,
  "type": "BUY"
}
```

**Response (400 Bad Request)** - Rejected ❌
```json
{
  "success": false,
  "error": {
    "code": "MAXIMUM_QUANTITY_EXCEEDED",
    "message": "maximum quantity exceeded: you can only buy/sell up to 10000 shares per transaction, but trying to BUY 10001",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 10001,
      "price": 100,
      "type": "BUY"
    }
  }
}
```

---

# STOCK AVAILABILITY LIMIT (Market Cap Protection)

## Overview
To prevent over-buying and ensure market integrity, each stock has a **total available quantity** in the market. Users cannot collectively buy more shares than what exists. When someone sells shares back, those shares become available for other users to buy.

---

## Test Case 14: Buy Above Available Stock in Market ❌

### Scenario
- Market has RELIANCE-CE-2900 with 100 total shares
- User 1 buys 80 shares (20 remain available)
- User 2 tries to buy 30 shares (only 20 available)

### Request by User 2
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user_2

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 30,
  "price": 2500,
  "type": "BUY"
}
```

### JSON Response (400 Bad Request) - INSUFFICIENT_STOCK
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "insufficient stock: only 20.00 shares of RELIANCE-CE-2900 available in market, but trying to buy 30.00",
    "details": {
      "symbol": "RELIANCE-CE-2900",
      "quantity": 30,
      "price": 2500,
      "type": "BUY"
    }
  }
}
```

### UI Display Suggestion
```
❌ Cannot Buy (Stock Unavailable)
Error Code: INSUFFICIENT_STOCK
Message: Only 20 shares available in market for RELIANCE-CE-2900
Requested: 30 shares
Available in Market: 20 shares
Shortfall: 10 shares

Action: You can:
- Buy 20 shares (all available)
- Wait for other traders to sell shares back into market
-Check if another similar stock is available
```

### Calculation Shown in Response
```
Total Shares in Market: 100
Already Held by Traders: 80 (can't be bought)
Available to Buy: 20
Your Request: 30
Status: INSUFFICIENT STOCK ❌
```

---

## Test Case 15: Buy Exactly Available Stock in Market ✅

### Scenario
- Market has RELIANCE-CE-2900 with 100 shares
- User buys exactly 100 shares (at market cap)

### Request
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 100,
  "price": 1000,
  "type": "BUY"
}
```

### JSON Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 100000,
    "fee": 500,
    "symbol": "RELIANCE-CE-2900",
    "quantity": 100,
    "price": 1000,
    "type": "BUY"
  }
}
```

### UI Display Suggestion
```
✅ Buy Successful (Market Cap Reached)
Symbol: RELIANCE-CE-2900
Quantity: 100 shares
Price: ₹1,000 per share
Total: ₹100,000
Fee: ₹500
Total Debit: ₹100,500

⚠️ Note: You now own 100% of market cap
No more shares available for purchase
```

---

## Test Case 16: Stock Becomes Available After Others Sell ✅

### Scenario - Step 1
- Market: 50 total shares of TCS-CE-3500
- User 1 buys all 50 shares
- User 2 tries to buy 1 share → FAILS (no stock available)

### Request by User 2 (Fails)
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user_2

{
  "symbol": "TCS-CE-3500",
  "quantity": 1,
  "price": 3500,
  "type": "BUY"
}
```

### Response (400 Bad Request) - No Stock Available
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "insufficient stock: only 0.00 shares of TCS-CE-3500 available in market, but trying to buy 1.00",
    "details": {
      "symbol": "TCS-CE-3500",
      "quantity": 1,
      "price": 3500,
      "type": "BUY"
    }
  }
}
```

### Scenario - Step 2: User 1 Sells 30 Shares
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user_1

{
  "symbol": "TCS-CE-3500",
  "quantity": 30,
  "price": 3600,
  "type": "SELL"
}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 108000,
    "fee": 540,
    "symbol": "TCS-CE-3500",
    "quantity": 30,
    "price": 3600,
    "type": "SELL"
  }
}
```

### Market State After Sell
```
Market Total: 50 shares
User 1 Holdings: 20 shares (50 - 30 sold)
Available in Market: 30 shares (just freed up)
```

### Scenario - Step 3: User 2 Buys (Now Succeeds!)
```json
POST /trading/trade
Content-Type: application/json
X-User-ID: test_user_2

{
  "symbol": "TCS-CE-3500",
  "quantity": 30,
  "price": 3500,
  "type": "BUY"
}
```

### Response (200 OK) - Now Available!
```json
{
  "success": true,
  "data": {
    "status": "Trade executed successfully",
    "total": 105000,
    "fee": 525,
    "symbol": "TCS-CE-3500",
    "quantity": 30,
    "price": 3500,
    "type": "BUY"
  }
}
```

### Final Market State
```
Market Total: 50 shares
User 1 Holdings: 20 shares
User 2 Holdings: 30 shares
Available in Market: 0 shares (all held)
```

### UI Flow Suggestion
```
Timeline View:
1️⃣  Initial State
    Market: 50 shares | In Holdings: 0 | Available: 50

2️⃣  User 1 Buys All
    Market: 50 shares | In Holdings: 50 | Available: 0 ❌

3️⃣  User 2 Tries to Buy
    ERROR: No stock available
    Message: "Only 0 shares available in market"

4️⃣  User 1 Sells 30
    Market: 50 shares | In Holdings: 20 | Available: 30 ✅

5️⃣  User 2 Buys 30
    SUCCESS! Stock became available
    User 2 Holdings: 30 shares
    Remaining Available: 0 shares
```

---

## Updated Error Codes Table (Including Stock Availability)


| Error Code | HTTP Status | Meaning | Limit |
|-----------|------------|---------|-------|
| `INSUFFICIENT_HOLDINGS` | 400 | Can't sell more than you own | Variable (by holdings) |
| `INSUFFICIENT_CASH` | 400 | Can't buy without funds | Variable (by cash available) |
| **`INSUFFICIENT_STOCK`** | **400** | **Can't buy more than market has** | **Variable (by market cap)** |
| `INVALID_QUANTITY_OR_PRICE` | 400 | Negative/zero values | qty > 0, price > 0 |
| **`MAXIMUM_QUANTITY_EXCEEDED`** | **400** | **Quantity exceeds limit** | **Max 10,000 per trade** |
| `INVALID_TRADE_TYPE` | 400 | Trade type not BUY/SELL | BUY or SELL only |
| `INVALID_REQUEST_BODY` | 400 | JSON parsing error | Valid JSON required |
| `METHOD_NOT_ALLOWED` | 405 | Wrong HTTP method | POST only |
| `UNAUTHORIZED` | 401 | Missing X-User-ID header | Header required |
| `WALLET_SNAPSHOT_ERROR` | 500 | Server error | Retry or contact support |
| `DEPOSIT_FAILED` | 400 | Deposit validation failed | amount > 0 |
| `INSUFFICIENT_FUNDS` | 400 | Insufficient balance for withdrawal | withdrawal ≤ available balance |


---

## Frontend UI Implementation for Quantity Limit

```javascript
// Validate quantity before sending request
const MAX_QUANTITY = 10000;

function validateQuantity(quantity) {
  const qty = parseFloat(quantity);
  
  if (qty <= 0) {
    return {
      valid: false,
      error: 'INVALID_QUANTITY_OR_PRICE',
      message: 'Quantity must be a positive number',
      suggestion: 'Enter a quantity greater than 0'
    };
  }
  
  if (qty > MAX_QUANTITY) {
    return {
      valid: false,
      error: 'MAXIMUM_QUANTITY_EXCEEDED',
      message: `Maximum ${MAX_QUANTITY.toLocaleString()} shares per transaction`,
      requested: qty.toLocaleString(),
      maximum: MAX_QUANTITY.toLocaleString(),
      suggestion: `You can split into ${Math.ceil(qty / MAX_QUANTITY)} transactions`
    };
  }
  
  return {
    valid: true,
    message: 'Quantity is valid'
  };
}

// Usage Example
const inputQty = document.getElementById('quantityInput').value;
const validation = validateQuantity(inputQty);

if (!validation.valid) {
  if (validation.error === 'MAXIMUM_QUANTITY_EXCEEDED') {
    showError(`Cannot ${tradeType.toUpperCase()}`);
    showDetail(`Requested: ${validation.requested} shares`);
    showDetail(`Maximum: ${validation.maximum} shares`);
    showSuggestion(`${validation.suggestion}`);
    displaySplitTransactionOption(inputQty, MAX_QUANTITY);
  }
} else {
  // Proceed with trade
  submitTrade(symbol, inputQty, price, tradeType);
}

// Helper function to show split transaction option
function displaySplitTransactionOption(totalQty, maxQty) {
  const numTransactions = Math.ceil(totalQty / maxQty);
  const remainder = totalQty % maxQty;
  
  let transactionList = '';
  for (let i = 1; i < numTransactions; i++) {
    transactionList += `Transaction ${i}: ${maxQty.toLocaleString()} shares\n`;
  }
  if (remainder > 0) {
    transactionList += `Transaction ${numTransactions}: ${remainder.toLocaleString()} shares`;
  } else {
    transactionList += `Transaction ${numTransactions}: ${maxQty.toLocaleString()} shares`;
  }
  
  console.log(`Suggested approach:\n${transactionList}`);
}

// Handle API response with MAXIMUM_QUANTITY_EXCEEDED error
fetch('/trading/trade', {...})
  .then(res => res.json())
  .then(data => {
    if (!data.success && data.error.code === 'MAXIMUM_QUANTITY_EXCEEDED') {
      const message = data.error.message;
      // Extract max allowed from error message
      const match = message.match(/up to (\d+)/);
      const maxAllowed = match ? parseInt(match[1]) : 10000;
      
      showUserError(
        'Quantity Limit Exceeded',
        `Maximum ${maxAllowed.toLocaleString()} shares allowed per transaction`,
        `Requested: ${data.error.details.quantity.toLocaleString()} shares`
      );
    }
  });
```

---

## Summary: Quantity Limits

- **Minimum**: > 0 shares (must be positive)
- **Maximum**: ≤ 10,000 shares per transaction
- **Special Cases**:
  - Exactly 10,000 shares: ✅ Allowed
  - 10,001+ shares: ❌ Rejected with `MAXIMUM_QUANTITY_EXCEEDED`
  - Different symbols don't share the limit (limit is per transaction, not per symbol)
  - Can execute multiple transactions to buy/sell more than 10,000 total

---

## Business Logic Behind 10,000 Share Limit

1. **Market Stability**: Prevents large single trades that could cause volatility
2. **Risk Management**: Limits exposure per transaction
3. **Settlement**: Easier to handle settlement and clearance
4. **Fair Trading**: Encourages distributed participation
5. **Platform Load**: Prevents overwhelming the trading engine

Users who need more shares should execute multiple transactions.

