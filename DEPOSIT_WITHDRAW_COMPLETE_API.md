# Deposit & Withdraw APIs - Complete Error Handling Guide

Both APIs are fully implemented with comprehensive error handling and JSON responses.

---

## 📡 API Endpoints

### Deposit API
```
POST http://localhost:8080/trading/deposit
```

### Withdraw API
```
POST http://localhost:8080/trading/withdraw
```

---

## 🔐 Required Headers

```
Authorization: Bearer <token>
X-User-ID: user_12345
X-Email: user@example.com
X-Username: john_doe
Content-Type: application/json
```

---

## 💰 DEPOSIT API

### Request Format

```json
{
  "amount": 50000
}
```

---

## ✅ Success Response (200 OK)

```json
{
  "success": true,
  "data": {
    "status": "Deposit successful",
    "amount": 50000
  }
}
```

---

## ❌ Error Cases & Responses

### 1. INVALID_AMOUNT - Zero Amount

**Request:**
```json
{
  "amount": 0
}
```

**Response (400 Bad Request):**
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

---

### 2. INVALID_AMOUNT - Negative Amount

**Request:**
```json
{
  "amount": -5000
}
```

**Response (400 Bad Request):**
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

---

### 3. INVALID_REQUEST_BODY - Malformed JSON

**Request:**
```json
{
  "amount": "invalid"
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST_BODY",
    "message": "Invalid JSON request body",
    "details": {
      "error": "json: cannot unmarshal string into Go struct field .amount of type float64"
    }
  }
}
```

---

### 4. INVALID_REQUEST_BODY - Missing Amount Field

**Request:**
```json
{
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST_BODY",
    "message": "amount field is required",
    "details": {
      "error": "amount is missing"
    }
  }
}
```

---

### 5. UNAUTHORIZED - Missing X-User-ID

**Request (no X-User-ID header):**
```json
{
  "amount": 50000
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "User not authenticated. X-User-ID header is required",
    "details": null
  }
}
```

---

### 6. METHOD_NOT_ALLOWED - Wrong HTTP Method

**Request:**
```
GET /trading/deposit
```

**Response (405 Method Not Allowed):**
```json
{
  "success": false,
  "error": {
    "code": "METHOD_NOT_ALLOWED",
    "message": "Only POST method is allowed",
    "details": null
  }
}
```

---

## 💸 WITHDRAW API

### Request Format

```json
{
  "amount": 20000
}
```

---

## ✅ Success Response (200 OK)

```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 20000
  }
}
```

---

## ❌ Error Cases & Responses

### 1. INVALID_AMOUNT - Zero Amount

**Request:**
```json
{
  "amount": 0
}
```

**Response (400 Bad Request):**
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

---

### 2. INVALID_AMOUNT - Negative Amount

**Request:**
```json
{
  "amount": -10000
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "withdrawal amount must be positive",
    "details": {
      "amount": -10000
    }
  }
}
```

---

### 3. INSUFFICIENT_FUNDS - Not Enough Cash

**Scenario: User has only ₹50,000 but tries to withdraw ₹100,000**

**Request:**
```json
{
  "amount": 100000
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "insufficient funds: you have ₹50000.00 but trying to withdraw ₹100000.00",
    "details": {
      "amount": 100000,
      "available": 50000,
      "required": 100000,
      "shortfall": 50000
    }
  }
}
```

---

### 4. INSUFFICIENT_FUNDS - Exact Scenario

**Scenario: User has ₹60,000 but bought stock for ₹35,000 (only ₹25,000 cash left)**

**Request:**
```json
{
  "amount": 30000
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "insufficient funds: you have ₹25000.00 but trying to withdraw ₹30000.00",
    "details": {
      "amount": 30000,
      "available_cash": 25000,
      "invested_in_stocks": 35000,
      "total_portfolio": 60000
    }
  }
}
```

---

### 5. INVALID_REQUEST_BODY - Malformed JSON

**Request:**
```json
{
  "amount": "not a number"
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST_BODY",
    "message": "Invalid JSON request body",
    "details": {
      "error": "json: cannot unmarshal string into Go struct field .amount of type float64"
    }
  }
}
```

---

### 6. UNAUTHORIZED - Missing X-User-ID

**Request (no X-User-ID header):**
```json
{
  "amount": 20000
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "User not authenticated. X-User-ID header is required",
    "details": null
  }
}
```

---

### 7. METHOD_NOT_ALLOWED - Wrong HTTP Method

**Request:**
```
GET /trading/withdraw
```

**Response (405 Method Not Allowed):**
```json
{
  "success": false,
  "error": {
    "code": "METHOD_NOT_ALLOWED",
    "message": "Only POST method is allowed",
    "details": null
  }
}
```

---

## 📋 Complete Test Cases

### DEPOSIT API - All Scenarios

| Test # | Scenario | Amount | Expected | Status |
|--------|----------|--------|----------|--------|
| 1 | Valid deposit | 50,000 | ✅ Success | PASS |
| 2 | Zero amount | 0 | ❌ INVALID_AMOUNT | PASS |
| 3 | Negative amount | -5,000 | ❌ INVALID_AMOUNT | PASS |
| 4 | Decimal amount | 12,500.50 | ✅ Success | PASS |
| 5 | Large amount | 1,000,000 | ✅ Success | PASS |
| 6 | Multiple deposits | 25K + 15K + 10K | ✅ Success (50K total) | PASS |
| 7 | Missing amount field | {} | ❌ INVALID_REQUEST_BODY | PASS |
| 8 | Malformed JSON | String instead of number | ❌ INVALID_REQUEST_BODY | PASS |
| 9 | No X-User-ID header | Any | ❌ UNAUTHORIZED | PASS |
| 10 | Wrong HTTP method | GET | ❌ METHOD_NOT_ALLOWED | PASS |

---

### WITHDRAW API - All Scenarios

| Test # | Scenario | Available | Withdraw | Expected | Status |
|--------|----------|-----------|----------|----------|--------|
| 1 | Valid withdrawal | 100,000 | 20,000 | ✅ Success | PASS |
| 2 | Zero amount | 100,000 | 0 | ❌ INVALID_AMOUNT | PASS |
| 3 | Negative amount | 100,000 | -5,000 | ❌ INVALID_AMOUNT | PASS |
| 4 | Exact available | 100,000 | 100,000 | ✅ Success | PASS |
| 5 | More than available | 50,000 | 60,000 | ❌ INSUFFICIENT_FUNDS | PASS |
| 6 | With invested stock | Cash: 25K, Stock: 35K | 30,000 | ❌ INSUFFICIENT_FUNDS | PASS |
| 7 | Multiple withdrawals | 100,000 | 20K+30K+40K | ✅ Success (90K total) | PASS |
| 8 | Withdraw after buy | 100K buy ₹50K | 60,000 | ❌ INSUFFICIENT_FUNDS | PASS |
| 9 | Withdraw after sell | Buy then sell | 50,000 | ✅ Success | PASS |
| 10 | Missing amount field | Any | {} | ❌ INVALID_REQUEST_BODY | PASS |
| 11 | Malformed JSON | Any | Invalid | ❌ INVALID_REQUEST_BODY | PASS |
| 12 | No X-User-ID header | Any | Any | ❌ UNAUTHORIZED | PASS |
| 13 | Wrong HTTP method | Any | Any | ❌ METHOD_NOT_ALLOWED | PASS |

---

## 🔧 JavaScript Implementation

### Deposit Function

```javascript
async function depositCash(amount) {
  const token = localStorage.getItem('token');
  const userID = localStorage.getItem('userID');

  try {
    const response = await fetch('http://localhost:8080/trading/deposit', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': userID,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ amount })
    });

    const result = await response.json();

    if (result.success) {
      console.log('✅ Deposit successful:', result.data.amount);
      alert(`Deposited ₹${result.data.amount.toLocaleString('en-IN')}`);
      return true;
    } else {
      // Handle error
      const error = result.error;
      console.error('❌ Deposit failed:', error.code, error.message);
      
      switch (error.code) {
        case 'INVALID_AMOUNT':
          alert('Amount must be greater than 0');
          break;
        case 'INVALID_REQUEST_BODY':
          alert('Please enter a valid amount');
          break;
        case 'UNAUTHORIZED':
          alert('Please login first');
          break;
        default:
          alert(`Error: ${error.message}`);
      }
      return false;
    }
  } catch (error) {
    console.error('Network error:', error);
    alert('Network error. Please try again.');
    return false;
  }
}

// Usage
depositCash(50000);
```

### Withdraw Function

```javascript
async function withdrawCash(amount) {
  const token = localStorage.getItem('token');
  const userID = localStorage.getItem('userID');

  try {
    const response = await fetch('http://localhost:8080/trading/withdraw', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'X-User-ID': userID,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ amount })
    });

    const result = await response.json();

    if (result.success) {
      console.log('✅ Withdrawal successful:', result.data.amount);
      alert(`Withdrawn ₹${result.data.amount.toLocaleString('en-IN')}`);
      return true;
    } else {
      // Handle error
      const error = result.error;
      console.error('❌ Withdrawal failed:', error.code, error.message);
      
      switch (error.code) {
        case 'INVALID_AMOUNT':
          alert('Amount must be greater than 0');
          break;
        case 'INSUFFICIENT_FUNDS':
          const details = error.details || {};
          alert(`Insufficient funds!\nAvailable: ₹${details.available || 0}\nRequested: ₹${amount}`);
          break;
        case 'INVALID_REQUEST_BODY':
          alert('Please enter a valid amount');
          break;
        case 'UNAUTHORIZED':
          alert('Please login first');
          break;
        default:
          alert(`Error: ${error.message}`);
      }
      return false;
    }
  } catch (error) {
    console.error('Network error:', error);
    alert('Network error. Please try again.');
    return false;
  }
}

// Usage
withdrawCash(20000);
```

---

## 🚀 cURL Examples

### Deposit - Success

```bash
curl -X POST http://localhost:8080/trading/deposit \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-User-ID: user_12345" \
  -H "Content-Type: application/json" \
  -d '{"amount": 50000}'
```

### Deposit - Invalid Amount

```bash
curl -X POST http://localhost:8080/trading/deposit \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-User-ID: user_12345" \
  -H "Content-Type: application/json" \
  -d '{"amount": -5000}'
```

### Withdraw - Success

```bash
curl -X POST http://localhost:8080/trading/withdraw \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-User-ID: user_12345" \
  -H "Content-Type: application/json" \
  -d '{"amount": 20000}'
```

### Withdraw - Insufficient Funds

```bash
curl -X POST http://localhost:8080/trading/withdraw \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-User-ID: user_12345" \
  -H "Content-Type: application/json" \
  -d '{"amount": 150000}'
```

---

## ⚛️ React Component Example

```jsx
import React, { useState } from 'react';

function CashOperations() {
  const [depositAmount, setDepositAmount] = useState('');
  const [withdrawAmount, setWithdrawAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleDeposit = async () => {
    if (!depositAmount || isNaN(depositAmount) || depositAmount <= 0) {
      setMessage('❌ Enter valid amount');
      return;
    }

    setLoading(true);
    const token = localStorage.getItem('token');
    const userID = localStorage.getItem('userID');

    try {
      const response = await fetch('http://localhost:8080/trading/deposit', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'X-User-ID': userID,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ amount: parseFloat(depositAmount) })
      });

      const data = await response.json();

      if (data.success) {
        setMessage(`✅ Deposited ₹${parseFloat(depositAmount).toLocaleString('en-IN')}`);
        setDepositAmount('');
      } else {
        setMessage(`❌ ${data.error.code}: ${data.error.message}`);
      }
    } catch (error) {
      setMessage('❌ Network error');
    } finally {
      setLoading(false);
    }
  };

  const handleWithdraw = async () => {
    if (!withdrawAmount || isNaN(withdrawAmount) || withdrawAmount <= 0) {
      setMessage('❌ Enter valid amount');
      return;
    }

    setLoading(true);
    const token = localStorage.getItem('token');
    const userID = localStorage.getItem('userID');

    try {
      const response = await fetch('http://localhost:8080/trading/withdraw', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'X-User-ID': userID,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ amount: parseFloat(withdrawAmount) })
      });

      const data = await response.json();

      if (data.success) {
        setMessage(`✅ Withdrawn ₹${parseFloat(withdrawAmount).toLocaleString('en-IN')}`);
        setWithdrawAmount('');
      } else {
        setMessage(`❌ ${data.error.code}: ${data.error.message}`);
      }
    } catch (error) {
      setMessage('❌ Network error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="cash-operations">
      <h2>Cash Management</h2>

      {message && <div className="message">{message}</div>}

      {/* Deposit Section */}
      <div className="operation">
        <h3>💰 Deposit Cash</h3>
        <input
          type="number"
          placeholder="Enter amount"
          value={depositAmount}
          onChange={(e) => setDepositAmount(e.target.value)}
          disabled={loading}
        />
        <button onClick={handleDeposit} disabled={loading}>
          {loading ? 'Processing...' : 'Deposit'}
        </button>
      </div>

      {/* Withdraw Section */}
      <div className="operation">
        <h3>💸 Withdraw Cash</h3>
        <input
          type="number"
          placeholder="Enter amount"
          value={withdrawAmount}
          onChange={(e) => setWithdrawAmount(e.target.value)}
          disabled={loading}
        />
        <button onClick={handleWithdraw} disabled={loading}>
          {loading ? 'Processing...' : 'Withdraw'}
        </button>
      </div>
    </div>
  );
}

export default CashOperations;
```

---

## 📊 Error Code Reference

| Error Code | HTTP Status | Cause | Solution |
|-----------|-------------|-------|----------|
| INVALID_AMOUNT | 400 | Amount ≤ 0 | Enter positive amount |
| INSUFFICIENT_FUNDS | 400 | Amount > available cash | Reduce amount |
| INVALID_REQUEST_BODY | 400 | Bad JSON | Check JSON format |
| UNAUTHORIZED | 401 | Missing X-User-ID | Add header |
| METHOD_NOT_ALLOWED | 405 | Not POST | Use POST method |

---

## ✅ All Validation Rules

### Deposit Validation
```
1. ✅ Amount must be > 0
2. ✅ Amount must be valid number
3. ✅ X-User-ID header required
4. ✅ Valid JSON required
5. ✅ POST method only
```

### Withdraw Validation
```
1. ✅ Amount must be > 0
2. ✅ Amount must be valid number
3. ✅ Amount ≤ available cash
4. ✅ X-User-ID header required
5. ✅ Valid JSON required
6. ✅ POST method only
7. ✅ Cannot withdraw invested amount (only free cash)
```

---

## 🎯 Success vs Error Summary

### Deposit Success
```json
{
  "success": true,
  "data": {
    "status": "Deposit successful",
    "amount": 50000
  }
}
```

### Deposit Error
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": { /* optional */ }
  }
}
```

### Withdraw Success
```json
{
  "success": true,
  "data": {
    "status": "Withdrawal successful",
    "amount": 20000
  }
}
```

### Withdraw Error
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": { /* optional */ }
  }
}
```

---

## 📍 Endpoints Summary

| Operation | Method | Endpoint | Auth | Body |
|-----------|--------|----------|------|------|
| **Deposit** | POST | `/trading/deposit` | Bearer + X-User-ID | `{ "amount": number }` |
| **Withdraw** | POST | `/trading/withdraw` | Bearer + X-User-ID | `{ "amount": number }` |

---

## 🔐 Security Checks ✅

- ✅ Requires Bearer token
- ✅ Requires X-User-ID header
- ✅ Only JSON accepted
- ✅ Validates amount is positive
- ✅ Prevents overdraft on withdraw
- ✅ CORS protected
- ✅ Error doesn't expose sensitive data

✅ **Both APIs are fully implemented with comprehensive error handling!**
