package websocket

import (
	"fmt"
	"log"
	"net/http"

	wsgorilla "github.com/gorilla/websocket"
)

type SClient struct {
	Connection *wsgorilla.Conn
	Socket     SSocket
}

func (this SClient) Active() bool {
	return this.Connection != nil
}

func (this *SClient) Start(port string) {
	if this.Active() {
		this.Printf("Can not start client while its already active.")
		return
	}

	this.Printf("---- Starting Server ----\n")

	http.HandleFunc("/", this.onConnect)
	this.Printf("server started on port %s\n", port)

	var err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		this.Printf("Error starting server listener: %v\n", err)
		return
	}

	this.Connection.SetCloseHandler(this.onClose)

	go this.sender()
	go this.listener()
}

func (this SClient) sender() {
	// this.Socket
}

func (this SClient) listener() {
	for {
		var _, message, err = this.Connection.ReadMessage()
		if err != nil {
			this.Printf("server: Error reading message: %v\n", err)
			break
		}

		this.Printf("Received message: %s\n", string(message))

		// OnResponse(message)
	}
}

func (this *SClient) onConnect(w http.ResponseWriter, r *http.Request) {
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

	this.Printf("Connected!\n")

	// TODO: Integrate this into client
	// for _, cb := range OnConnectionCallbacks {
	// 	cb(this.Connection)
	// }
}

func (this *SClient) onClose(code int, text string) error {
	this.Printf("onClose: %d - %s\n", code, text)
	this.Connection = nil
	return nil
}

func (this SClient) Printf(format string, args ...any) {
	format = fmt.Sprintf("client: %s", format)
	log.Printf(format, args...)
}
