package websocket

import (
	"net/http"
	"sync"

	wsgorilla "github.com/gorilla/websocket"
)

type SClient struct {
	Connection *wsgorilla.Conn
	Socket     SSocket

	onReadyNotify []*chan bool
	mutex         sync.Mutex
}

func (this *SClient) Active() bool {
	return this.Connection != nil
}

func (this *SClient) Start(port string) {
	if this.Active() {
		clog.Errorf("Can not start client while its already active.")
		return
	}

	clog.Infof("---- Starting Server ----\n")

	http.HandleFunc("/", this.onConnect)
	clog.Infof("server started on port %s\n", port)

	var err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		clog.Errorf("Error starting server listener: %v\n", err)
		return
	}
}

func (this *SClient) Close() {
	if !this.Active() {
		clog.Error("Unable to close the connection while it is nil\n")
		return
	}

	if err := this.Connection.CloseHandler()(0, "exit"); err != nil {
		clog.Errorf("Error while attempting to close the conntection: %v\n", err)
		this.Connection.Close()
	}
}

func (this *SClient) sender() {
	for {
		var message = <-this.Socket.Channel

		if this.Socket.Channel == nil {
			break
		}

		this.Connection.WriteJSON(message.Request)
	}
}

func (this *SClient) listener() {
	for {
		var _, message, err = this.Connection.ReadMessage()
		if err != nil {
			clog.Errorf("Error reading message: %v\n", err)
			break
		}

		this.Socket.receive(message)
	}
}

func (this *SClient) onReady() {
	this.mutex.Lock()
	clog.Infof("Ready!")

	this.Connection.SetCloseHandler(this.onClose)
	this.Socket.Open()

	if len(this.onReadyNotify) > 0 {
		// Unblock any scripts waiting on this signal
		for _, sub := range this.onReadyNotify {
			*sub <- true
		}
	}

	go this.sender()
	go this.listener()
	this.mutex.Unlock()
}

func (this *SClient) OnReadySub() *chan bool {
	this.mutex.Lock()
	if this.onReadyNotify == nil {
		this.onReadyNotify = make([]*chan bool, 0)
	}

	var sub = make(chan bool)
	this.onReadyNotify = append(this.onReadyNotify, &sub)

	this.mutex.Unlock()

	return &sub
}

func (this *SClient) onConnect(w http.ResponseWriter, r *http.Request) {
	if this.Active() {
		clog.Infof("Overwriting existing connections with new one\n")
	}

	var upgrader = wsgorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Accept any connection
			return true
		},
	}

	var err error
	// Upgrade the HTTP connection to a WebSocket connection
	this.Connection, err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		clog.Fatalf("Error upgrading: %v\n", err)
		return
	}

	this.onReady()
}

func (this *SClient) onClose(code int, text string) error {
	clog.Infof("onClose: %d - %s\n", code, text)
	this.Connection = nil
	this.Socket.Close()
	return nil
}
