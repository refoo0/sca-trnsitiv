package websocketserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket" // WebSocket library
)

// WebSocketServer struct holds the server configuration
type WebSocketServer struct {
	Port     string
	Upgrader websocket.Upgrader
}

// NewWebSocketServer creates a new instance of the WebSocketServer with default configuration
func NewWebSocketServer(port string) *WebSocketServer {
	return &WebSocketServer{
		Port: port,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Security issue: This allows connections from any source (Cross-Origin Request Forgery, CSRF)
				// This should be modified to allow connections only from specific origins
				return true
			},
		},
	}
}

// EchoHandler handles incoming WebSocket connections and echoes messages back to the client
func (s *WebSocketServer) EchoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		fmt.Printf("Message received: %s\n", message)

		// Send message back to client (echo)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

// Start starts the WebSocket server and listens for incoming connections
func (s *WebSocketServer) Start() error {
	fmt.Printf("Starting WebSocket server on port %s...\n", s.Port)

	http.HandleFunc("/ws", s.EchoHandler)

	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
	return err
}
