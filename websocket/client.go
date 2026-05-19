package websocket

import (
	"fmt"
	"log"
	"net/http"

	wsgorilla "github.com/gorilla/websocket"
)

type SClient struct {
	Connection *wsgorilla.Conn
}

func (this SClient) Active() bool {
	return this.Connection != nil
}

func (this *SClient) Start(port string) {
	if this.Active() {
		this.Printf("Can not start client while its already active.")
		return
	}

	StartServer(port, this.connectionHandler)
}
func (this SClient) listener() {
	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := this.Connection.ReadMessage()
		if err != nil {
			this.Printf("server: Error reading message: %v\n", err)
			break
		}

		this.Printf("Received message: %s\n", string(message))

		// OnResponse(message)
	}
}

func (this *SClient) connectionHandler(w http.ResponseWriter, r *http.Request) {
	if this.Active() {
		this.Printf("Overwriting existing connections with new one\n")
	}

	var upgrader = wsgorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// accept any connection
			return true
		},
	}

	var err error
	// Upgrade the HTTP connection to a WebSocket connection
	this.Connection, err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		this.Printf("Error upgrading: %v\n", err)
		return
	}

	defer this.Connection.Close()

	this.Printf("Connected!\n")

	defer func() { this.Connection = nil }()

	// TODO: Integrate this into client
	for _, cb := range OnConnectionCallbacks {
		cb(this.Connection)
	}

	this.listener()
}

func (this SClient) Printf(format string, args ...any) {
	format = fmt.Sprintf("client: %s", format)
	log.Printf(format, args...)
}
