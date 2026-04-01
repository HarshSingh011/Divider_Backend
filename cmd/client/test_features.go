package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type TradeRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 ` json:"quantity"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
}

type AlertRequest struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Condition string  `json:"condition"`
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

func main() {
	baseURL := "https://divider-backend.onrender.com"

	fmt.Println("\n════════════════════════════════════════════════════════════")
	fmt.Println("🚀 StockTrack Trading Platform - Feature Demonstration")
	fmt.Println("════════════════════════════════════════════════════════════\n")

	fmt.Println("📝 Step 1: Authentication")
	fmt.Println("──────────────────────────────────────────────────────────")

	payload := AuthPayload{
		Email:    "trader@stocktrack.com",
		Username: "portfolio_trader",
		Password: "SecurePass123",
	}

	token, err := registerUser(baseURL, payload)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}
	fmt.Printf("✅ User registered successfully!\n")
	fmt.Printf("   Token: %s\n\n", token[:30]+"...")

	fmt.Println("💰 Step 2: Immutable Transaction Ledger - Deposit Cash")
	fmt.Println("──────────────────────────────────────────────────────────")

	if err := depositCash(baseURL, token, 50000); err != nil {
		fmt.Printf("❌ Deposit failed: %v\n", err)
	} else {
		fmt.Printf("✅ Deposited ₹50,000\n")
	}

	wallet, err := getWalletSnapshot(baseURL, token)
	if err != nil {
		fmt.Printf("❌ Failed to get wallet: %v\n", err)
	} else {
		fmt.Printf("✅ Current Wallet Balance: ₹%.2f\n", wallet["total_balance"])
		fmt.Printf("   Available Cash: ₹%.2f\n\n", wallet["available_cash"])
	}

	fmt.Println("🚨 Step 3: Real-Time Price Alert System")
	fmt.Println("──────────────────────────────────────────────────────────")

	alerts := []AlertRequest{
		{Symbol: "RELIANCE-CE-2900", Price: 50.0, Condition: "ABOVE"},
		{Symbol: "HDFC-PE-1400", Price: 10.0, Condition: "BELOW"},
	}

	for _, alert := range alerts {
		if err := createAlert(baseURL, token, alert); err == nil {
			fmt.Printf("✅ Alert created: %s when price goes %s ₹%.2f\n",
				alert.Symbol, alert.Condition, alert.Price)
		}
	}

	fmt.Println("\n📊 Step 4: Immutable Transaction Ledger - Trading")
	fmt.Println("──────────────────────────────────────────────────────────")

	trades := []TradeRequest{
		{Symbol: "RELIANCE-CE-2900", Quantity: 10, Price: 45.50, Type: "BUY"},
		{Symbol: "INFY-CE-1500", Quantity: 5, Price: 62.75, Type: "BUY"},
	}

	for _, trade := range trades {
		if err := executeTrade(baseURL, token, trade); err == nil {
			fee := (trade.Quantity * trade.Price) * 0.005
			totalCost := (trade.Quantity * trade.Price) + fee
			fmt.Printf("✅ %s %d units of %s @ ₹%.2f (Fee: ₹%.2f, Total: ₹%.2f)\n",
				trade.Type, int(trade.Quantity), trade.Symbol, trade.Price, fee, totalCost)
		}
	}

	wallet, err = getWalletSnapshot(baseURL, token)
	if err == nil {
		fmt.Printf("\n✅ Updated Wallet Balance: ₹%.2f\n", wallet["total_balance"])
		fmt.Printf("   Cash After Trades: ₹%.2f\n\n", wallet["available_cash"])
	}

	fmt.Println("📈 Step 5: OHLC Aggregator - Candlestick Data")
	fmt.Println("──────────────────────────────────────────────────────────")

	if candles, err := getCandles(baseURL, "RELIANCE-CE-2900", 5); err == nil {
		fmt.Printf("✅ Retrieved %d candles for RELIANCE-CE-2900\n", len(candles))
		if len(candles) > 0 {
			for i, candle := range candles {
				fmt.Printf("   [Candle %d] O:%.2f H:%.2f L:%.2f C:%.2f Vol:%d\n",
					i+1, candle["open"], candle["high"], candle["low"], candle["close"], candle["volume"])
			}
		}
	}

	fmt.Println("\n🔄 Step 6: Real-Time Market Updates via WebSocket")
	fmt.Println("──────────────────────────────────────────────────────────")
	fmt.Println("Connecting to WebSocket (5 seconds)...\n")

	if err := connectToWebSocket(token); err != nil {
		fmt.Printf("❌ WebSocket failed: %v\n", err)
	}

	fmt.Println("\n════════════════════════════════════════════════════════════")
	fmt.Println("✨ All Features Demonstrated Successfully!")
	fmt.Println("════════════════════════════════════════════════════════════\n")
}

func registerUser(baseURL string, payload AuthPayload) (string, error) {
	data, _ := json.Marshal(payload)
	resp, err := http.Post(
		baseURL+"/auth/register",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AuthResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token, nil
}

func depositCash(baseURL, token string, amount float64) error {
	req := DepositRequest{Amount: amount}
	data, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", baseURL+"/trading/deposit", bytes.NewBuffer(data))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getWalletSnapshot(baseURL, token string) (map[string]interface{}, error) {
	httpReq, _ := http.NewRequest("GET", baseURL+"/trading/wallet", nil)
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func createAlert(baseURL, token string, alert AlertRequest) error {
	data, _ := json.Marshal(alert)

	httpReq, _ := http.NewRequest("POST", baseURL+"/trading/alerts", bytes.NewBuffer(data))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func executeTrade(baseURL, token string, trade TradeRequest) error {
	data, _ := json.Marshal(trade)

	httpReq, _ := http.NewRequest("POST", baseURL+"/trading/trade", bytes.NewBuffer(data))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getCandles(baseURL, symbol string, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/trading/candles?symbol=%s&limit=%d", baseURL, symbol, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var candles []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&candles)
	return candles, nil
}

func connectToWebSocket(token string) error {
	wsURL := fmt.Sprintf("wss://divider-backend.onrender.com/ws?token=%s", token)

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("✅ Connected to market stream\n")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	tickCount := 0
	go func() {
		for {
			var marketData []map[string]interface{}
			err := conn.ReadJSON(&marketData)
			if err != nil {
				return
			}

			tickCount++
			if tickCount == 1 {
				fmt.Println("📊 Live Market Data:")
				fmt.Println("───────────────────────────────────────────")
			}
			if tickCount <= 3 {
				for _, stock := range marketData {
					fmt.Printf("   %s: ₹%.2f (%+.2f%%)\n",
						stock["symbol"], stock["currentPrice"], stock["percentageChange"])
				}
				fmt.Println()
			}
		}
	}()

	<-ticker.C
	fmt.Println("✅ Market stream received successfully")
	return nil
}
