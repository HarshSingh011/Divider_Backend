package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"stocktrack-backend/internal/domain"
)

// Hub manages all WebSocket connections
type Hub struct {
	// Map of connected clients
	clients map[*Client]bool
	// Broadcast channel for market data
	broadcast chan []domain.MarketTick
	// Channel for client registrations
	register chan *Client
	// Channel for client unregistrations
	unregister chan *Client
	// Mutex to protect clients map
	mu sync.Mutex
}

// Client represents a connected WebSocket client
type Client struct {
	conn     *websocket.Conn
	hub      *Hub
	send     chan []domain.MarketTick
	userID   string
	email    string
	username string
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []domain.MarketTick, 10),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start runs the hub's event loop
func (h *Hub) Start() {
	go func() {
		for {
			select {
			// New client connection
			case client := <-h.register:
				h.mu.Lock()
				h.clients[client] = true
				h.mu.Unlock()
				fmt.Printf("[WS] Client connected: %s (%s)\n", client.userID, client.email)
				fmt.Printf("[WS] Total clients: %d\n", len(h.clients))

			// Client disconnection
			case client := <-h.unregister:
				h.mu.Lock()
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
				h.mu.Unlock()
				fmt.Printf("[WS] Client disconnected: %s\n", client.userID)
				fmt.Printf("[WS] Total clients: %d\n", len(h.clients))

			// Broadcast market data to all clients
			case marketData := <-h.broadcast:
				h.mu.Lock()
				for client := range h.clients {
					select {
					case client.send <- marketData:
						// Sent successfully
					default:
						// Client's send channel is full, skip
					}
				}
				h.mu.Unlock()
			}
		}
	}()
}

// Stop gracefully shuts down the hub
func (h *Hub) Stop() {
	h.mu.Lock()
	for client := range h.clients {
		delete(h.clients, client)
		close(client.send)
	}
	h.mu.Unlock()
}

// Broadcast sends market data to all connected clients
func (h *Hub) Broadcast(marketData []domain.MarketTick) {
	h.broadcast <- marketData
}

// Handler is the HTTP handler for WebSocket connections
func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	// Extract user info from headers (set by middleware)
	userID := r.Header.Get("X-User-ID")
	email := r.Header.Get("X-Email")
	username := r.Header.Get("X-Username")

	if userID == "" {
		http.Error(w, "Missing user information", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[WS] Upgrade error: %v\n", err)
		return
	}

	// Create new client
	client := &Client{
		conn:     conn,
		hub:      h,
		send:     make(chan []domain.MarketTick, 10),
		userID:   userID,
		email:    email,
		username: username,
	}

	// Register client with hub
	h.register <- client

	// Start reading and writing goroutines
	go client.readLoop()
	go client.writeLoop()
}

// readLoop reads messages from the client (for keep-alive)
func (c *Client) readLoop() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var tick domain.MarketTick
		err := c.conn.ReadJSON(&tick)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("[WS] Error: %v\n", err)
			}
			return
		}
	}
}

// writeLoop sends messages to the client
func (c *Client) writeLoop() {
	for {
		select {
		case marketData, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(marketData)
			if err != nil {
				fmt.Printf("[WS] Write error: %v\n", err)
				return
			}
		}
	}
}
