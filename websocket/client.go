package websocket

import (
	"net/http"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	wsgorilla "github.com/gorilla/websocket"
)

var clog = clogger.SClog{
	Prefix: "client: ",
	LogState: map[clogger.ELogLevel]func() bool{
		clogger.LogAll:                     func() bool { return true },
		clogger.LogInfo | clogger.LogError: func() bool { return true },
	},
}

type SClient struct {
	Connection *wsgorilla.Conn
	Socket     SSocket
	Ready      chan struct{}
}

func (this SClient) Active() bool {
	return this.Connection != nil
}

func (this *SClient) Start(port string) {
	if this.Active() {
		clog.Printf("Can not start client while its already active.")
		return
	}

	clog.Printf("---- Starting Server ----\n")

	http.HandleFunc("/", this.onConnect)
	clog.Printf("server started on port %s\n", port)

	var err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		clog.Printf("Error starting server listener: %v\n", err)
		return
	}
}

func (this SClient) Close() {
	if !this.Active() {
		clog.Printf("Unable to close the connection while it is nil\n")
		return
	}

	this.Connection.Close()
}

func (this *SClient) sender() {
	for {
		var message = <-this.Socket.Channel
		// clog.Printf("sending message: %v\n", message)
		this.Connection.WriteJSON(message.Request)
	}
}

func (this *SClient) listener() {
	for {
		var _, message, err = this.Connection.ReadMessage()
		if err != nil {
			clog.Printf("server: Error reading message: %v\n", err)
			break
		}

		// clog.Printf("Received message: %s\n", string(message))
		this.Socket.Receive(message)
	}
}

func (this *SClient) onReady() {
	clog.Printf("Ready!")
	this.Connection.SetCloseHandler(this.onClose)

	this.Socket.Open()

	if this.Ready != nil {
		// unblock any scripts waiting on this signal
		close(this.Ready)
	}

	go this.sender()
	go this.listener()
}

func (this *SClient) onConnect(w http.ResponseWriter, r *http.Request) {
	if this.Active() {
		clog.Printf("Overwriting existing connections with new one\n")
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
		clog.Printf("Error upgrading: %v\n", err)
		return
	}

	// TODO: Integrate this into client
	// for _, cb := range OnConnectionCallbacks {
	// 	cb(this.Connection)
	// }

	this.onReady()
}

func (this *SClient) onClose(code int, text string) error {
	clog.Printf("onClose: %d - %s\n", code, text)
	this.Connection = nil
	this.Socket.Close()
	return nil
}

// func (this SClient) Printf(format string, args ...any) {
// 	format = fmt.Sprintf("client: %s", format)
// 	log.Printf(format, args...)
// }
