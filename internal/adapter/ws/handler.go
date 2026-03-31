package ws

import (
	"fmt"
	"net/http"
	"sync"

	"stocktrack-backend/internal/domain"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan []domain.MarketTick
	register chan *Client
	unregister chan *Client
	mu sync.Mutex
}

type Client struct {
	conn     *websocket.Conn
	hub      *Hub
	send     chan []domain.MarketTick
	userID   string
	email    string
	username string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []domain.MarketTick, 10),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Start() {
	go func() {
		for {
			select {
			case client := <-h.register:
				h.mu.Lock()
				h.clients[client] = true
				h.mu.Unlock()
				fmt.Printf("[WS] Client connected: %s (%s)\n", client.userID, client.email)
				fmt.Printf("[WS] Total clients: %d\n", len(h.clients))

			case client := <-h.unregister:
				h.mu.Lock()
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
				h.mu.Unlock()
				fmt.Printf("[WS] Client disconnected: %s\n", client.userID)
				fmt.Printf("[WS] Total clients: %d\n", len(h.clients))

			case marketData := <-h.broadcast:
				h.mu.Lock()
				for client := range h.clients {
					select {
					case client.send <- marketData:
					default:
					}
				}
				h.mu.Unlock()
			}
		}
	}()
}

func (h *Hub) Stop() {
	h.mu.Lock()
	for client := range h.clients {
		delete(h.clients, client)
		close(client.send)
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(marketData []domain.MarketTick) {
	h.broadcast <- marketData
}

func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	email := r.Header.Get("X-Email")
	username := r.Header.Get("X-Username")

	if userID == "" {
		http.Error(w, "Missing user information", http.StatusUnauthorized)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[WS] Upgrade error: %v\n", err)
		return
	}

	client := &Client{
		conn:     conn,
		hub:      h,
		send:     make(chan []domain.MarketTick, 10),
		userID:   userID,
		email:    email,
		username: username,
	}

	h.register <- client

	go client.readLoop()
	go client.writeLoop()
}

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
