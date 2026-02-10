package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
	"github.com/sidji-omnichannel/internal/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// Connect upgrades HTTP connection to WebSocket
// GET /api/ws
func (h *WebSocketHandler) Connect(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orgIDValue, _ := c.Get("organization_id")

	userID := userIDValue.(uuid.UUID)
	orgID := orgIDValue.(uuid.UUID)

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Create client
	client := &websocket.Client{
		Hub:            h.hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		OrganizationID: orgID,
	}

	// Register client
	h.hub.Register(client)

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump()
}
