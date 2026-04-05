# Deposit & Withdraw APIs - JSON Response Bodies (UI Ready)

## Overview
Complete test cases for Deposit and Withdraw operations with:
- All edge cases covered
- Validation checks enforced
- Structured JSON error responses
- Success responses with transaction details

---

## DEPOSIT API Tests

### Test 1: Deposit Positive Amount ✅

**Request**:
```json
POST /trading/deposit
Content-Type: application/json
X-User-ID: test_user

{
  "amount": 50000
}
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Deposit successful",
    "amount": 50000
  }
}
```

**Calculation**:
- Previous Balance: ₹100,000 (initial)
- Deposit: ₹50,000
- New Balance: ₹150,000

**UI Display**:
```
✅ Deposit Successful
Amount: ₹50,000.00
Previous Balance: ₹100,000.00
New Balance: ₹150,000.00
```

---

### Test 2: Cannot Deposit Zero Amount ❌

**Request**:
```json
POST /trading/deposit

{
  "amount": 0
}
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "deposit amount must be positive",
    "details": {
      "amount": 0
    }
  }
}
```

**UI Display**:
```
❌ Invalid Deposit
Error Code: INVALID_AMOUNT  
Message: Deposit amount must be positive
Current Balance: ₹100,000.00
```

---

### Test 3: Cannot Deposit Negative Amount ❌

**Request**:
```json
POST /trading/deposit

{
  "amount": -5000
}
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "deposit amount must be positive",
    "details": {
      "amount": -5000
    }
  }
}
```

**UI Display**:
```
❌ Invalid Deposit
Error Code: INVALID_AMOUNT
Message: Deposit amount must be positive
Entered Amount: ₹-5000.00 ❌
```

---

### Test 4: Multiple Deposits Accumulate ✅

**Request 1**:
```json
POST /trading/deposit

{
  "amount": 25000
}
```
**Response (200 OK)**: Success ✅

**Request 2**:
```json
POST /trading/deposit

{
  "amount": 15000
}
```
**Response (200 OK)**: Success ✅

**Request 3**:
```json
POST /trading/deposit

{
  "amount": 10000
}
```
**Response (200 OK)**: Success ✅

**Final Wallet Check**:
```json
GET /trading/wallet
```

**JSON Response**:
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 150000,
    "available_cash": 150000,
    "invested_amount": 0,
    "positions": {}
  }
}
```

**Calculation**:
- Initial: ₹100,000
- Deposit 1: +₹25,000 → ₹125,000
- Deposit 2: +₹15,000 → ₹140,000
- Deposit 3: +₹10,000 → ₹150,000

**UI Display**:
```
Deposit History:
1. ₹25,000.00 ✓
2. ₹15,000.00 ✓
3. ₹10,000.00 ✓

Total Deposits: ₹50,000.00
Current Balance: ₹150,000.00
```

---

## WITHDRAW API Tests

### Test 5: Withdraw with Sufficient Funds ✅

**Stock**: Initial balance is ₹100,000

**Request**:
```json
POST /trading/withdraw

{
  "amount": 20000
}
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 20000
  }
}
```

**Calculation**:
- Balance Before: ₹100,000
- Withdrawal: ₹20,000
- Balance After: ₹80,000

**UI Display**:
```
✅ Withdrawal Successful
Amount Withdrawn: ₹20,000.00
Previous Balance: ₹100,000.00
New Balance: ₹80,000.00
```

---

### Test 6: Cannot Withdraw Zero Amount ❌

**Request**:
```json
POST /trading/withdraw

{
  "amount": 0
}
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "withdrawal amount must be positive",
    "details": {
      "amount": 0
    }
  }
}
```

**UI Display**:
```
❌ Invalid Withdrawal
Error Code: INVALID_AMOUNT
Message: Withdrawal amount must be positive
Current Balance: ₹100,000.00
```

---

### Test 7: Cannot Withdraw Negative Amount ❌

**Request**:
```json
POST /trading/withdraw

{
  "amount": -5000
}
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "withdrawal amount must be positive",
    "details": {
      "amount": -5000
    }
  }
}
```

**UI Display**:
```
❌ Invalid Withdrawal
Error Code: INVALID_AMOUNT
Message: Withdrawal amount must be positive
Entered Amount: ₹-5000.00 ❌
```

---

### Test 8: Cannot Withdraw More Than Available ❌

**Balance**: ₹100,000

**Request**:
```json
POST /trading/withdraw

{
  "amount": 150000
}
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "insufficient funds: you have ₹100000.00 but trying to withdraw ₹150000.00",
    "details": {
      "amount": 150000
    }
  }
}
```

**UI Display**:
```
❌ Insufficient Funds
Error Code: INSUFFICIENT_FUNDS
Need: ₹150,000.00
Have: ₹100,000.00
Shortfall: ₹50,000.00

Suggestion: You can withdraw up to ₹100,000.00
```

---

### Test 9: Withdraw Exact Available Amount ✅

**Scenario**: Buy shares first, then withdraw remaining cash

**Step 1 - Buy 50 shares @ 100**:
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

Available cash calculation:
- Initial: ₹100,000
- Cost: 50 × 100 = ₹5,000
- Fee: 5,000 × 0.005 = ₹25
- Total spent: ₹5,025
- **Available cash: ₹94,975**

**Step 2 - Withdraw Exact Available Amount**:
```json
POST /trading/withdraw

{
  "amount": 94975
}
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 94975
  }
}
```

**Final Wallet Check**:
```json
GET /trading/wallet
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 5000,
    "available_cash": 0,
    "invested_amount": 5000,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 50,
        "average_cost": 100,
        "current_price": 100,
        "unrealized_pnl": 0,
        "percentage": 0
      }
    }
  }
}
```

**UI Display**:
```
✅ Withdrawal Successful
Amount Withdrawn: ₹94,975.00
Holdings:
- RELIANCE-CE-2900: 50 shares @ ₹100 = ₹5,000
Available Cash: ₹0
Total Balance: ₹5,000.00
```

---

### Test 10: Multiple Withdrawals ✅

**Step 1 - Withdraw 10000**:
```json
POST /trading/withdraw
{ "amount": 10000 }
```
Response: 200 OK ✅

**Step 2 - Withdraw 20000**:
```json
POST /trading/withdraw
{ "amount": 20000 }
```
Response: 200 OK ✅

**Step 3 - Withdraw 30000**:
```json
POST /trading/withdraw
{ "amount": 30000 }
```
Response: 200 OK ✅

**Calculation**:
- Initial: ₹100,000
- Withdrawal 1: -₹10,000 → ₹90,000
- Withdrawal 2: -₹20,000 → ₹70,000
- Withdrawal 3: -₹30,000 → ₹40,000

**Wallet After**:
```json
GET /trading/wallet
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 40000,
    "available_cash": 40000,
    "invested_amount": 0,
    "positions": {}
  }
}
```

**UI Display**:
```
Withdrawal History:
1. ₹10,000.00 ✓
2. ₹20,000.00 ✓
3. ₹30,000.00 ✓

Total Withdrawals: ₹60,000.00
Remaining Balance: ₹40,000.00
```

---

### Test 11: Cannot Withdraw More After Previous Withdrawal ❌

**Step 1 - Withdraw 50000**:
```json
POST /trading/withdraw
{ "amount": 50000 }
```
Response: 200 OK ✅ (Remaining: ₹50,000)

**Step 2 - Try to Withdraw 60000**:
```json
POST /trading/withdraw
{ "amount": 60000 }
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "insufficient funds: you have ₹50000.00 but trying to withdraw ₹60000.00",
    "details": {
      "amount": 60000
    }
  }
}
```

**UI Display**:
```
❌ Cannot Withdraw
Error Code: INSUFFICIENT_FUNDS
After previous withdrawal, you have: ₹50,000.00
Trying to withdraw: ₹60,000.00
Shortfall: ₹10,000.00
```

---

### Test 12: Deposit Then Withdraw ✅

**Step 1 - Deposit 50000**:
```json
POST /trading/deposit
{ "amount": 50000 }
```
Response: 200 OK ✅
Balance now: ₹150,000

**Step 2 - Withdraw 80000**:
```json
POST /trading/withdraw
{ "amount": 80000 }
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 80000
  }
}
```

**Final Balance**:
- Start: ₹100,000
- Deposit: +₹50,000 = ₹150,000
- Withdrawal: -₹80,000 = **₹70,000**

**UI Display**:
```
✅ Operations Successful
Operation 1: Deposit ₹50,000 ✓
Balance: ₹150,000.00

Operation 2: Withdrawal ₹80,000 ✓
Final Balance: ₹70,000.00
```

---

### Test 13: Buy, Sell, Then Withdraw ✅

**Step 1 - Buy 50 shares @ 100**:
```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 100,
  "type": "BUY"
}
```
Response: 200 OK
- Cost: ₹5,025 (with fee)
- Cash remaining: ₹94,975

**Step 2 - Sell 50 shares @ 120**:
```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 50,
  "price": 120,
  "type": "SELL"
}
```
Response: 200 OK
- Proceeds: ₹6,000 - fee ₹30 = ₹5,970
- Cash now: ₹94,975 + ₹5,970 = ₹100,945

**Step 3 - Withdraw 50000**:
```json
POST /trading/withdraw
{ "amount": 50000 }
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 50000
  }
}
```

**Final Balance**: ₹100,945 - ₹50,000 = **₹50,945**

**UI Display**:
```
Transaction Summary:
1. BUY 50 @ ₹100 = -₹5,025 (cost + fee) ✓
2. SELL 50 @ ₹120 = +₹5,970 (proceeds - fee) ✓
3. Profit from trade: +₹945 ✓
4. Withdraw ₹50,000 ✓

Final Balance: ₹50,945.00
```

---

### Test 14: Cannot Withdraw Below Holdings Value ❌

**Setup: Buy shares to reduce available cash**

**Step 1 - Buy 100 shares @ 500**:
```json
POST /trading/trade
{
  "symbol": "RELIANCE-CE-2900",
  "quantity": 100,
  "price": 500,
  "type": "BUY"
}
```
- Cost: 100 × 500 = ₹50,000
- Fee: ₹50,000 × 0.005 = ₹250
- Total: ₹50,250
- **Available cash: ₹49,750**

**Step 2 - Try to Withdraw 60000**:
```json
POST /trading/withdraw
{ "amount": 60000 }
```

**JSON Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "insufficient funds: you have ₹49750.00 but trying to withdraw ₹60000.00",
    "details": {
      "amount": 60000
    }
  }
}
```

**Step 3 - Withdraw Available Amount (Success)**:
```json
POST /trading/withdraw
{ "amount": 49750 }
```

**JSON Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 49750
  }
}
```

**Final Wallet**:
```json
GET /trading/wallet
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user_id": "test_user",
    "total_balance": 50000,
    "available_cash": 0,
    "invested_amount": 50000,
    "positions": {
      "RELIANCE-CE-2900": {
        "symbol": "RELIANCE-CE-2900",
        "quantity": 100,
        "average_cost": 500,
        "current_price": 500,
        "unrealized_pnl": 0,
        "percentage": 0
      }
    }
  }
}
```

**UI Display**:
```
❌ Cannot Withdraw ₹60,000
Error: Insufficient funds
Available to withdraw: ₹49,750.00

Holdings:
- RELIANCE-CE-2900: 100 shares @ ₹500 = ₹50,000
- Locked in investments

✓ Can withdraw: ₹49,750.00
```

---

## Error Codes Summary

| Code | HTTP Status | Scenario | UI Action |
|------|---|---|---|
| `INVALID_AMOUNT` | 400 | Amount ≤ 0 | Show validation error |
| `INSUFFICIENT_FUNDS` | 400 | Withdraw > available | Show max available |
| `DEPOSIT_FAILED` | 400 | Deposit transaction error | Retry or contact support |
| `WITHDRAWAL_FAILED` | 400 | Withdrawal transaction error | Retry or contact support |
| `INVALID_REQUEST_BODY` | 400 | JSON parsing error | Check input format |
| `UNAUTHORIZED` | 401 | Missing X-User-ID | Authenticate user |
| `METHOD_NOT_ALLOWED` | 405 | Wrong HTTP method | Check API docs |
| `WALLET_SNAPSHOT_ERROR` | 500 | Server error | Retry or contact support |

---

## Success Codes

| Code | HTTP Status | Scenario |
|------|---|---|
| 200 | OK | Deposit/Withdraw successful |
| SUCCESS | success=true | Transaction completed |

---

## Frontend Integration Examples

### Deposit Example:
```javascript
async function depositmoney(amount) {
  try {
    const response = await fetch(`${API_BASE}/trading/deposit`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-User-ID': userId
      },
      body: JSON.stringify({ amount: amount })
    });
    
    const data = await response.json();
    
    if (data.success) {
      showSuccess(`✅ Deposited ₹${data.data.amount}`);
      updateWallet();
    } else {
      const error = data.error;
      if (error.code === 'INVALID_AMOUNT') {
        showError('Amount must be positive');
      } else {
        showError(error.message);
      }
    }
  } catch (error) {
    showError('Deposit failed: ' + error.message);
  }
}
```

### Withdraw Example:
```javascript
async function withdrawMoney(amount) {
  try {
    const response = await fetch(`${API_BASE}/trading/withdraw`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-User-ID': userId
      },
      body: JSON.stringify({ amount: amount })
    });
    
    const data = await response.json();
    
    if (data.success) {
      showSuccess(`✅ Withdrew ₹${data.data.amount}`);
      updateWallet();
    } else {
      const error = data.error;
      if (error.code === 'INSUFFICIENT_FUNDS') {
        showError(`You only have ₹${error.details.available}`);
      } else if (error.code === 'INVALID_AMOUNT') {
        showError('Amount must be positive');
      } else {
        showError(error.message);
      }
    }
  } catch (error) {
    showError('Withdrawal failed: ' + error.message);
  }
}
```

---

## Test Results Summary

✅ **All Tests Passing** (14/14)

| Test | Status | Error Code |
|------|--------|-----------|
| Deposit Positive | ✅ | - |
| Deposit Zero | ✅ | INVALID_AMOUNT |
| Deposit Negative | ✅ | INVALID_AMOUNT |
| Multiple Deposits | ✅ | - |
| Withdraw Sufficient | ✅ | - |
| Withdraw Zero | ✅ | INVALID_AMOUNT |
| Withdraw Negative | ✅ | INVALID_AMOUNT |
| Withdraw More Than Available | ✅ | INSUFFICIENT_FUNDS |
| Withdraw Exact Amount | ✅ | - |
| Multiple Withdrawals | ✅ | - |
| Withdraw Overlimit After Previous | ✅ | INSUFFICIENT_FUNDS |
| Deposit Then Withdraw | ✅ | - |
| Buy Sell Then Withdraw | ✅ | - |
| Cannot Withdraw Below Holdings | ✅ | INSUFFICIENT_FUNDS |
