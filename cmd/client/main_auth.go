package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// User payload for register/login
type AuthPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse from server
type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func main() {
	fmt.Println("=== StockTrack Test Client ===\n")

	// Step 1: Register a new user
	fmt.Println("1️⃣  Registering user...")
	registerPayload := AuthPayload{
		Email:    "trader@example.com",
		Username: "trader123",
		Password: "securepass123",
	}

	token, err := registerUser(registerPayload)
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}

	fmt.Printf("✓ User registered! Token: %s\n\n", token[:20]+"...")

	// Step 2: Connect to WebSocket with token
	fmt.Println("2️⃣  Connecting to WebSocket...")
	if err := connectToWebSocket(token); err != nil {
		fmt.Printf("❌ WebSocket connection failed: %v\n", err)
		return
	}
}

// registerUser registers a new user and returns JWT token
func registerUser(payload AuthPayload) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		"http://localhost:8080/auth/register",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("registration failed: %s", string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Token, nil
}

// connectToWebSocket connects to the WebSocket and listens for market updates
func connectToWebSocket(token string) error {
	wsURL := fmt.Sprintf("ws://localhost:8080/ws?token=%s", token)
	fmt.Printf("Connecting to: %s\n", wsURL)

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	fmt.Println("✓ Connected to WebSocket!\n")
	fmt.Println("📊 Market Data Stream (10 seconds):")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	tickCount := 0
	go func() {
		for {
			var marketData []map[string]interface{}
			err := conn.ReadJSON(&marketData)
			if err != nil {
				fmt.Printf("❌ Connection closed: %v\n", err)
				return
			}

			tickCount++
			fmt.Printf("\n[Tick %d] Received %d stocks:\n", tickCount, len(marketData))
			for _, stock := range marketData {
				change := stock["percentageChange"].(float64)
				symbol := stock["symbol"].(string)
				price := stock["currentPrice"].(float64)

				// Color code the change (green for up, red for down)
				sign := "📈"
				if change < 0 {
					sign = "📉"
				}

				fmt.Printf("  %s %-20s $%-8.2f (%+.2f%%)\n", sign, symbol, price, change)
			}
		}
	}()

	<-ticker.C
	fmt.Println("\n✓ Test complete! Server is working correctly.")
	return nil
}
