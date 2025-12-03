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

func StartServer(port string) {
	print("\n---- Starting Server ----\n")

	http.HandleFunc("/", wsHandler)
	fmt.Printf("server: Started on: %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("server: Error starting server:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("server: Error upgrading:", err)
		return
	}

	defer conn.Close()

	fmt.Print("server: Connected to client!\n")

	if ActiveConnection != nil {
		fmt.Printf("server: Overwriting existing connections with new one\n")
	}

	ActiveConnection = conn
	defer func() { ActiveConnection = nil }()

	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("server: Error reading message:", err)
			break
		}

		OnResponse(message)
	}
}
