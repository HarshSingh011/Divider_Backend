#!/usr/bin/env python3
"""
Test script to verify sell validation:
1. Cannot sell more shares than you own
2. Cannot sell shares you don't own at all
3. Can successfully sell the exact amount you own
"""

import requests
import json
import time

BASE_URL = "http://localhost:5000"
USER_ID = "test_user_123"
SYMBOL = "RELIANCE-CE-2900"

def set_auth_headers():
    """Create auth headers with dummy token"""
    return {
        "Content-Type": "application/json",
        "Authorization": "Bearer test_token_123",
        "X-User-ID": USER_ID
    }

def buy_shares(quantity, price):
    """Buy shares"""
    print(f"\n📈 BUYING {quantity} shares of {SYMBOL} @ ₹{price}")
    response = requests.post(
        f"{BASE_URL}/trading/trade",
        headers=set_auth_headers(),
        json={
            "symbol": SYMBOL,
            "quantity": quantity,
            "price": price,
            "type": "BUY"
        }
    )
    print(f"Status: {response.status_code}")
    if response.status_code != 200:
        print(f"❌ Error: {response.text}")
        return False
    print(f"✅ Success")
    return True

def sell_shares(quantity, price):
    """Attempt to sell shares"""
    print(f"\n📉 SELLING {quantity} shares of {SYMBOL} @ ₹{price}")
    response = requests.post(
        f"{BASE_URL}/trading/trade",
        headers=set_auth_headers(),
        json={
            "symbol": SYMBOL,
            "quantity": quantity,
            "price": price,
            "type": "SELL"
        }
    )
    print(f"Status: {response.status_code}")
    if response.status_code != 200:
        print(f"❌ Error: {response.text}")
        return False
    print(f"✅ Success")
    return True

def get_holdings():
    """Get current wallet snapshot"""
    print(f"\n💼 CHECKING WALLET")
    response = requests.get(
        f"{BASE_URL}/trading/wallet",
        headers=set_auth_headers()
    )
    if response.status_code == 200:
        data = response.json()
        positions = data.get("positions", {})
        if SYMBOL in positions:
            qty = positions[SYMBOL].get("quantity", 0)
            print(f"Holdings: {qty} shares of {SYMBOL}")
            print(f"Full data: {json.dumps(positions[SYMBOL], indent=2)}")
        else:
            print(f"No holdings of {SYMBOL}")
        print(f"Available cash: ₹{data.get('available_cash', 0):.2f}")
    else:
        print(f"❌ Error: {response.text}")

def test_scenario():
    print("=" * 60)
    print("SELL VALIDATION TEST SUITE")
    print("=" * 60)
    
    # Test 1: Cannot sell what you don't own
    print("\n\n🧪 TEST 1: Try to sell shares you don't own")
    print("-" * 50)
    get_holdings()
    result1 = sell_shares(50, 2700)
    
    if result1:
        print("⚠️  BUG FOUND! System allowed selling shares you don't own!")
    else:
        print("✅ Correctly rejected (expected)")
    
    # Test 2: Buy shares
    print("\n\n🧪 TEST 2: Buy 100 shares")
    print("-" * 50)
    buy_shares(100, 2500)
    get_holdings()
    
    # Test 3: Sell exact amount
    print("\n\n🧪 TEST 3: Sell exact amount owned (100 shares)")
    print("-" * 50)
    result3 = sell_shares(100, 2700)
    get_holdings()
    
    if result3:
        print("✅ Correctly allowed (expected)")
    else:
        print("❌ Incorrectly rejected")
    
    # Test 4: Try to sell more than owned
    print("\n\n🧪 TEST 4: Buy 50 shares then try to sell 75")
    print("-" * 50)
    buy_shares(50, 2600)
    get_holdings()
    result4 = sell_shares(75, 2700)
    
    if result4:
        print("⚠️  BUG FOUND! System allowed overselling!")
        get_holdings()
    else:
        print("✅ Correctly rejected (expected)")
    
    # Test 5: Sell remaining
    print("\n\n🧪 TEST 5: Sell remaining shares (should succeed)")
    print("-" * 50)
    result5 = sell_shares(50, 2700)
    get_holdings()
    
    if result5:
        print("✅ Correctly allowed (expected)")
    else:
        print("❌ Incorrectly rejected")

if __name__ == "__main__":
    print("Waiting 2 seconds for server to be ready...")
    time.sleep(2)
    try:
        test_scenario()
    except Exception as e:
        print(f"\n❌ Test error: {e}")
        import traceback
        traceback.print_exc()
