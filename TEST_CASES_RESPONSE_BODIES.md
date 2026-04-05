# Stock Trading API - 9 Test Cases with Response Bodies

## Overview
This document shows all 9 sell validation test cases with:
- The request being made
- The HTTP response status code
- The response body

---

## Test Case 1: Sell Non-Owned Shares ❌
**Scenario**: User tries to sell shares they don't own at all

### Request
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 2700,
  "type": "SELL"
}
```

### Response
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient holdings: you have 0.00 shares of RELIANCE-CE-2900 but trying to sell 10.00
```

### Status Code: 400 Bad Request
**Error Message**: User has NO holdings of this stock, so the sell is rejected

---

## Test Case 2: Sell More Than Owned ❌
**Scenario**: User has 10 shares but tries to sell 15

### Request (Step 1: Buy 10 shares first)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY",
  "total": 8000
}
```

### Request (Step 2: Try to sell 15)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 15,
  "price": 900,
  "type": "SELL"
}
```

### Response (Step 2: Rejected)
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient holdings: you have 10.00 shares of RELIANCE-CE-2900 but trying to sell 15.00
```

### Status Code: 400 Bad Request
**Error Message**: Trying to oversell - have 10 but trying to sell 15

---

## Test Case 3: Sell Exact Amount Owned ✅
**Scenario**: User has 10 shares and sells exactly 10 shares

### Request (Step 1: Buy 10 shares)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 800,
  "type": "BUY",
  "total": 8000
}
```

### Request (Step 2: Sell all 10 shares)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 900,
  "type": "SELL"
}
```

### Response (Step 2: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 900,
  "type": "SELL",
  "total": 9000
}
```

### Status Code: 200 OK
**Success**: User can sell exactly the amount they own (10 = 10) ✅

---

## Test Case 4: Sell Partial Amount ✅
**Scenario**: User has 100 shares but sells only 60

### Request (Step 1: Buy 100 shares @ 100)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 100,
  "price": 100,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 100,
  "price": 100,
  "type": "BUY",
  "total": 10000
}
```

### Request (Step 2: Sell 60 shares @ 120)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 60,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 2: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 60,
  "price": 120,
  "type": "SELL",
  "total": 7200
}
```

### Wallet After (Step 3: Check holdings)
```http
GET /trading/wallet HTTP/1.1
X-User-ID: test_user
```

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
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
```

### Status Code: 200 OK
**Success**: User can sell partial amount (60 out of 100) and retains 40 shares ✅

---

## Test Case 5: Cannot Sell After Holdings Empty ❌
**Scenario**: User sells all shares, then tries to sell again

### Request (Step 1: Buy 10 shares)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 100,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 100,
  "type": "BUY",
  "total": 1000
}
```

### Request (Step 2: Sell all 10 shares)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 2: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 10,
  "price": 120,
  "type": "SELL",
  "total": 1200
}
```

### Request (Step 3: Try to sell again)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 5,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 3: Rejected)
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient holdings: you have 0.00 shares of RELIANCE-CE-2900 but trying to sell 5.00
```

### Status Code: 400 Bad Request
**Error Message**: Holdings are now empty (0 shares), cannot sell any more ❌

---

## Test Case 6: Buy Insufficient Cash ❌
**Scenario**: User tries to buy 1000 shares @ 200 but only has 100,000 initial cash

### Request
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 1000,
  "price": 200,
  "type": "BUY"
}
```

### Response
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient cash: you have ₹100000.00 but need ₹201000.00
```

**Calculation**:
- Stock cost: 1000 × 200 = ₹200,000
- Fee (0.5%): 200,000 × 0.005 = ₹1,000
- Total needed: ₹201,000
- Available: ₹100,000
- **Result: Cannot buy** ❌

### Status Code: 400 Bad Request

---

## Test Case 7: Multiple Buys Then Sell ✅
**Scenario**: User buys shares multiple times, then sells

### Request (Step 1: Buy 30 shares @ 100)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 30,
  "price": 100,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 30,
  "price": 100,
  "type": "BUY",
  "total": 3000
}
```

### Request (Step 2: Buy 20 shares @ 110)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 20,
  "price": 110,
  "type": "BUY"
}
```

### Response (Step 2: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 20,
  "price": 110,
  "type": "BUY",
  "total": 2200
}
```

### Request (Step 3: Buy 50 shares @ 90)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 90,
  "type": "BUY"
}
```

### Response (Step 3: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 90,
  "type": "BUY",
  "total": 4500
}
```

### Request (Step 4: Try to sell 150 - should FAIL)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 150,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 4: Rejected)
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient holdings: you have 100.00 shares of RELIANCE-CE-2900 but trying to sell 150.00
```

### Request (Step 5: Sell 75 - should SUCCEED)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 75,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 5: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 75,
  "price": 120,
  "type": "SELL",
  "total": 9000
}
```

### Wallet After (Step 6: Check holdings)
```http
GET /trading/wallet HTTP/1.1
X-User-ID: test_user
```

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
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
```

### Status Code: 200 OK (for successful sell)
**Success**: Total holdings = 100 (30+20+50). Can't sell 150, but can sell 75. Remaining = 25 ✅

---

## Test Case 8: Multiple Symbols Independent ✅
**Scenario**: Trading different stocks doesn't affect other holdings

### Request (Step 1: Buy RELIANCE)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 100,
  "type": "BUY"
}
```

### Response (Step 1: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 100,
  "type": "BUY",
  "total": 5000
}
```

### Request (Step 2: Buy TCS)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "TCS-CE-3500",
  "quantity": 30,
  "price": 150,
  "type": "BUY"
}
```

### Response (Step 2: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "TCS-CE-3500",
  "quantity": 30,
  "price": 150,
  "type": "BUY",
  "total": 4500
}
```

### Request (Step 3: Sell 40 RELIANCE - should succeed)
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 40,
  "price": 120,
  "type": "SELL"
}
```

### Response (Step 3: Success)
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "Trade executed successfully",
  "symbol": "RELIANCE-CE-2900",
  "quantity": 40,
  "price": 120,
  "type": "SELL",
  "total": 4800
}
```

### Request (Step 4: Try to sell 50 TCS - should FAIL (only have 30))
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "TCS-CE-3500",
  "quantity": 50,
  "price": 170,
  "type": "SELL"
}
```

### Response (Step 4: Rejected)
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

insufficient holdings: you have 30.00 shares of TCS-CE-3500 but trying to sell 50.00
```

### Wallet After (Step 5: Check holdings)
```http
GET /trading/wallet HTTP/1.1
X-User-ID: test_user
```

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
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
```

### Status Code: 400 Bad Request (for oversell attempt)
**Success**: RELIANCE holdings = 10, TCS holdings = 30 (tracked independently) ✅

---

## Test Case 9: Invalid Quantity or Price ❌
**Scenario**: Attempting to trade with invalid values

### Request 1: Zero Quantity
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 0,
  "price": 100,
  "type": "BUY"
}
```

### Response 1: Rejected
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

quantity and price must be positive
```

---

### Request 2: Negative Quantity
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": -50,
  "price": 100,
  "type": "BUY"
}
```

### Response 2: Rejected
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

quantity and price must be positive
```

---

### Request 3: Zero Price
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 0,
  "type": "BUY"
}
```

### Response 3: Rejected
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

quantity and price must be positive
```

---

### Request 4: Negative Price
```http
POST /trading/trade HTTP/1.1
Content-Type: application/json
X-User-ID: test_user

{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": -100,
  "type": "BUY"
}
```

### Response 4: Rejected
```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

quantity and price must be positive
```

### Status Code: 400 Bad Request (for all invalid inputs)
**Result**: All invalid inputs are rejected ✅

---

## Summary Table

| Test Case | Type | Request | Response | Status |
|-----------|------|---------|----------|--------|
| 1 | Sell non-owned | SELL 10 (have 0) | Rejected | 400 ❌ |
| 2 | Sell oversell | SELL 15 (have 10) | Rejected | 400 ❌ |
| 3 | Sell exact | SELL 10 (have 10) | Success | 200 ✅ |
| 4 | Sell partial | SELL 60 (have 100) | Success | 200 ✅ |
| 5 | Sell empty | SELL 5 (have 0) | Rejected | 400 ❌ |
| 6 | Buy insufficient cash | BUY 1000 @ 200 | Rejected | 400 ❌ |
| 7 | Multiple buys | SELL 75 (have 100) | Success | 200 ✅ |
| 8 | Multi-symbol | SELL 50 TCS (have 30) | Rejected | 400 ❌ |
| 9 | Invalid inputs | Various invalid | Rejected | 400 ❌ |

---

## Key Points

✅ **Success (200 OK)**: Tests 3, 4, 7 return successful response with transaction details
❌ **Error (400 Bad Request)**: Tests 1, 2, 5, 6, 8, 9 return validation error messages
📊 **Holdings Tracking**: Response includes exact quantities and P&L calculations
💰 **Fee Calculation**: 0.5% fee is included in all transactions automatically
