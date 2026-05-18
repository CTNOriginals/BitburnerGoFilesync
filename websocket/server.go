package websocket

import (
	"log"
	"net/http"

	gorillaws "github.com/gorilla/websocket"
)

var OnConnectionCallbacks []func(ws *gorillaws.Conn) = []func(ws *gorillaws.Conn){}

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ActiveConnection *gorillaws.Conn

func StartServer(port string) {
	log.Print("\n---- Starting Server ----\n")

	http.HandleFunc("/", wsHandler)
	log.Printf("server: Started on: %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("server: Error starting server:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("server: Error upgrading:", err)
		return
	}

	defer conn.Close()

	log.Print("server: Connected to client!\n")

	if ActiveConnection != nil {
		log.Printf("server: Overwriting existing connections with new one\n")
	}

	ActiveConnection = conn
	defer func() { ActiveConnection = nil }()

	for _, cb := range OnConnectionCallbacks {
		cb(conn)
	}

	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("server: Error reading message:", err)
			break
		}

		OnResponse(message)
	}
}
