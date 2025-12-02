package communication

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ActiveConnection *websocket.Conn

func StartServer(port int) {
	print("\n---- Starting Server ----\n")

	http.HandleFunc("/", wsHandler)
	fmt.Printf("WebSocket server started on: %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	if ActiveConnection != nil {
		fmt.Printf("Overwriting existing connections with new one\n")
	}

	ActiveConnection = conn
	defer func() { ActiveConnection = nil }()

	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// fmt.Printf("\n🔻 Received: %s\n\n", message)

		OnResponse(message)
	}
}
