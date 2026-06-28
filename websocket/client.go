package websocket

import (
	"net/http"
	"slices"
	"sync"

	wsgorilla "github.com/gorilla/websocket"
)

type FnOnReadyCallback func() (unsub bool)

type SClient struct {
	Connection *wsgorilla.Conn
	Socket     *SSocket

	onReadyNotify   []*chan bool
	onReadyCallback []FnOnReadyCallback

	mutex sync.Mutex
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
	defer this.mutex.Unlock()

	if this.Socket == nil {
		this.Socket = &SSocket{}
	}

	this.Connection.SetCloseHandler(this.onClose)
	this.Socket.Open()

	// Unblock any scripts waiting on this signal
	for i, sub := range this.onReadyNotify {
		*sub <- true
		close(*sub)
		this.onReadyNotify = slices.Delete(this.onReadyNotify, i, i+1)
	}

	// Call each subscribed callback
	for i, sub := range this.onReadyCallback {
		// if callback returns true: unsubscribe
		if !sub() {
			this.onReadyCallback = slices.Delete(this.onReadyCallback, i, i+1)
		}
	}

	go this.sender()
	go this.listener()

	clog.Infof("Ready!")
}

func (this *SClient) OnReadySubNotify() *chan bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.onReadyNotify == nil {
		this.onReadyNotify = make([]*chan bool, 0)
	}

	var sub = make(chan bool)
	this.onReadyNotify = append(this.onReadyNotify, &sub)

	return &sub
}

// Subscribe with a callback that will be called each time [SClient.onReady] is called.
func (this *SClient) OnReadySubCallback(cb FnOnReadyCallback) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.onReadyCallback == nil {
		this.onReadyCallback = make([]FnOnReadyCallback, 0)
	}

	this.onReadyCallback = append(this.onReadyCallback, cb)
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
