package websocket

import (
	"log"
	"net/http"

	gorillaws "github.com/gorilla/websocket"
)

type TConnectionCallback func(ws *gorillaws.Conn)
type TConnectionCallbackList []TConnectionCallback

func (this *TConnectionCallbackList) Push(cb TConnectionCallback) int {
	var index = len(*this)

	*this = append(*this, cb)

	return index
}

func (this *TConnectionCallbackList) Remove(index int) {
	if index >= len(*this) {
		return
	}

	*this = append((*this)[0:index], (*this)[index+1:]...)
}

var OnConnectionCallbacks TConnectionCallbackList = make(TConnectionCallbackList, 0)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ActiveConnection *gorillaws.Conn

func StartServer(port string, handler func(http.ResponseWriter, *http.Request)) {
	log.Print("\n---- Starting Server ----\n")

	http.HandleFunc("/", handler)
	log.Printf("server: Started on: %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("server: Error starting server:", err)
	}
}

// Default websocket handler
func DefaultWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = gorillaws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// accept any connection
			return true
		},
	}
	// Upgrade the HTTP connection to a WebSocket connection
	var conn, err = upgrader.Upgrade(w, r, nil)

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

		log.Printf("Received message: %s\n", string(message))

		// OnResponse(message)
	}
}
