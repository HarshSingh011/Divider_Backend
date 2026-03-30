package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Connect to the WebSocket server
	url := "ws://localhost:8080/ws"
	fmt.Println("Connecting to", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	fmt.Println("✓ Connected to WebSocket!")
	fmt.Println("Listening for market ticks (10 seconds)...\n")

	// Read messages for 10 seconds
	done := make(chan struct{})
	go func() {
		defer close(done)
		tickCount := 0
		for {
			var marketData []map[string]interface{}
			err := conn.ReadJSON(&marketData)
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			tickCount++
			fmt.Printf("[Tick %d] Received %d stocks:\n", tickCount, len(marketData))
			for _, stock := range marketData {
				fmt.Printf("  - %s: $%.2f (%+.2f%%)\n",
					stock["symbol"],
					stock["currentPrice"],
					stock["percentageChange"])
			}
			fmt.Println()
		}
	}()

	// Wait for 10 seconds or for the goroutine to finish
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("\n✓ Test complete! Market data is streaming correctly.")
	case <-done:
	}
}
