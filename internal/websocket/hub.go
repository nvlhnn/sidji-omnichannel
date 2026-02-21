package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 10240
)

// Client represents a connected WebSocket client
type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by organization
	clients map[uuid.UUID]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast message to organization
	broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	OrganizationID uuid.UUID
	Event          string
	Data           interface{}
}

// WebSocketMessage represents the structure of WebSocket messages
type WebSocketMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.OrganizationID] == nil {
				h.clients[client.OrganizationID] = make(map[*Client]bool)
			}
			h.clients[client.OrganizationID][client] = true
			h.mu.Unlock()
			log.Printf("Client connected: user=%s org=%s", client.UserID, client.OrganizationID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.OrganizationID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.OrganizationID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: user=%s org=%s", client.UserID, client.OrganizationID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.OrganizationID]
			h.mu.RUnlock()

			msg, err := json.Marshal(map[string]interface{}{
				"event": message.Event,
				"data":  message.Data,
			})
			if err != nil {
				continue
			}

			for client := range clients {
				select {
				case client.Send <- msg:
				default:
					h.mu.Lock()
					close(client.Send)
					delete(h.clients[client.OrganizationID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// Broadcast sends a message to all clients in an organization
func (h *Hub) Broadcast(orgID uuid.UUID, event string, data interface{}) {
	h.broadcast <- &BroadcastMessage{
		OrganizationID: orgID,
		Event:          event,
		Data:           data,
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error (user=%s): %v", c.UserID, err)
			} else {
				log.Printf("WebSocket closed (user=%s): %v", c.UserID, err)
			}
			break
		}

		// Handle incoming messages (if needed)
		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Process client messages (e.g., typing indicators)
		switch msg.Event {
		case "typing":
			// Broadcast typing indicator to org
			c.Hub.Broadcast(c.OrganizationID, "typing", map[string]interface{}{
				"user_id": c.UserID,
				"data":    msg.Data,
			})
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
