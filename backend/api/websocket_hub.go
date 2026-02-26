//go:build web

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// wsMaxMessageSize limits incoming WebSocket messages to 1MB
	wsMaxMessageSize = 1 << 20
	// wsWriteWait is the time allowed to write a message
	wsWriteWait = 10 * time.Second
	// wsMaxClients limits concurrent WebSocket connections
	wsMaxClients = 50
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// WSHub manages all WebSocket connections
type WSHub struct {
	clients map[*websocket.Conn]bool
	mutex   sync.RWMutex
}

var hub *WSHub
var onceHub sync.Once

func Hub() *WSHub {
	if hub == nil {
		onceHub.Do(func() {
			hub = &WSHub{
				clients: make(map[*websocket.Conn]bool),
			}
		})
	}
	return hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Origin validation is handled by wsAuthCheck middleware
		// Allow all here to avoid double-checking
		return true
	},
}

// Emit sends an event to all connected WebSocket clients
func (h *WSHub) Emit(event string, data any) {
	msg := WSMessage{Event: event, Data: data}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for conn := range h.clients {
		conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Printf("ws write error: %v", err)
		}
	}
}

// HandleWebSocket handles WebSocket upgrade and connection lifecycle
func (h *WSHub) HandleWebSocket(c *gin.Context) {
	// Check max clients
	h.mutex.RLock()
	clientCount := len(h.clients)
	h.mutex.RUnlock()
	if clientCount >= wsMaxClients {
		c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "too many connections"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	// Set read limits to prevent oversized messages
	conn.SetReadLimit(wsMaxMessageSize)

	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		delete(h.clients, conn)
		h.mutex.Unlock()
		conn.Close()
	}()

	// read loop - handle incoming messages (e.g. CLI input)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		h.handleIncoming(msg)
	}
}

// handleIncoming processes messages from the client
func (h *WSHub) handleIncoming(msg WSMessage) {
	// dispatch CLI input events etc.
	if handler, ok := incomingHandlers[msg.Event]; ok {
		handler(msg.Data)
	}
}

var incomingHandlers = map[string]func(data any){}

// RegisterHandler registers a handler for incoming WebSocket events
func RegisterHandler(event string, handler func(data any)) {
	incomingHandlers[event] = handler
}
