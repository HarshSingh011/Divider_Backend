package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MarketTick represents a single price update for a stock or option
type MarketTick struct {
	Symbol           string    `json:"symbol"`
	CurrentPrice     float64   `json:"currentPrice"`
	PercentageChange float64   `json:"percentageChange"`
	Timestamp        time.Time `json:"timestamp"`
}

// Hub manages all active WebSocket connections and broadcasts market data
type Hub struct {
	// A map holding all active WebSocket connections
	clients map[*websocket.Conn]bool
	// A channel that receives new market data to broadcast
	broadcast chan []MarketTick
	// Channel to handle client registrations
	register chan *websocket.Conn
	// Channel to handle client unregistrations
	unregister chan *websocket.Conn
	// Mutex to lock the clients map to prevent race conditions
	mu sync.Mutex
}

// NewHub creates and returns a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []MarketTick),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// RunMarketEngine simulates the stock market by generating random price changes
func (h *Hub) RunMarketEngine() {
	// Create a mock list of stocks
	stocks := []MarketTick{
		{Symbol: "RELIANCE-CE-2900", CurrentPrice: 45.50, PercentageChange: 0},
		{Symbol: "HDFC-PE-1400", CurrentPrice: 12.20, PercentageChange: 0},
		{Symbol: "INFY-CE-1500", CurrentPrice: 62.75, PercentageChange: 0},
		{Symbol: "TCS-PE-3500", CurrentPrice: 85.30, PercentageChange: 0},
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Ticker fires every 500 milliseconds
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		<-ticker.C // Wait for the tick

		// Randomly change the prices slightly
		for i := range stocks {
			oldPrice := stocks[i].CurrentPrice
			change := (rand.Float64() - 0.5) * 2 // Random change between -1.0 and +1.0
			stocks[i].CurrentPrice += change
			stocks[i].PercentageChange = ((stocks[i].CurrentPrice - oldPrice) / oldPrice) * 100
			stocks[i].Timestamp = time.Now()
		}

		// Send the updated array into the broadcast channel
		h.broadcast <- stocks
	}
}

// StartBroadcasting listens for market data and sends it to all connected clients
func (h *Hub) StartBroadcasting() {
	for {
		select {
		// New client connection
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			fmt.Printf("[CONNECTED] Total clients: %d\n", len(h.clients))

		// Client disconnection
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			fmt.Printf("[DISCONNECTED] Total clients: %d\n", len(h.clients))

		// New market data to broadcast
		case marketData := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				// Send the data as JSON to the client
				err := client.WriteJSON(marketData)
				if err != nil {
					// If sending fails (e.g., user closed the app), remove them
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	// Allow connections from any origin (important for development)
	CheckOrigin: func(r *http.Request) bool { return true },
}

// serveWs handles WebSocket connections and upgrades HTTP to WS
func serveWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}

	// Register the new connection with the hub
	h.register <- conn

	// Listen for messages from the client (to detect disconnections)
	go func() {
		defer func() {
			h.unregister <- conn
		}()

		for {
			var tick MarketTick
			err := conn.ReadJSON(&tick)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("WebSocket error:", err)
				}
				return
			}
		}
	}()
}

// HealthCheck endpoint for debugging
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "Server is running", "time": "%s"}`, time.Now().String())
}

func main() {
	hub := NewHub()

	// Start the background workers (Goroutines)
	go hub.RunMarketEngine()
	go hub.StartBroadcasting()

	// Set up the API routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/health", healthCheck)

	// Start the server
	port := ":8080"
	fmt.Printf("\n=== Market Server running on %s ===\n", port)
	fmt.Println("WebSocket endpoint: ws://localhost:8080/ws")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("\nWaiting for connections...\n")

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
