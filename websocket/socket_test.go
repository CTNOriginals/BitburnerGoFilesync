package websocket

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOpenClose(t *testing.T) {
	var socket SSocket

	socket.Open()
	if !socket.isOpen {
		t.Error("isOpen = false after Open()")
	}
	if socket.Channel == nil {
		t.Error("Channel is nil after Open()")
	}
	if socket.Messages == nil {
		t.Error("Messages is nil after Open()")
	}

	socket.Close()
	if socket.isOpen {
		t.Error("isOpen = true after Close()")
	}
	if socket.Channel != nil {
		t.Error("Channel is not nil after Close()")
	}
}

func TestOpenCloseGuards(t *testing.T) {
	var socket SSocket

	socket.Open()
	socket.Open() // should log error, not panic

	socket.Close()
	socket.Close() // should log error, not panic
}

func TestGetId(t *testing.T) {
	var socket SSocket
	socket.Open()
	defer socket.Close()

	var id0 = socket.getId()
	var id1 = socket.getId()
	var id2 = socket.getId()

	if id0 != 0 {
		t.Errorf("first id = %d, want 0", id0)
	}
	if id1 != 1 {
		t.Errorf("second id = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("third id = %d, want 2", id2)
	}
}

func TestSend(t *testing.T) {
	var socket SSocket
	socket.Open()
	t.Cleanup(socket.Close)

	go func() { <-socket.Channel }()

	var msg = socket.send(PushFile, nil)
	if msg == nil {
		t.Fatal("send returned nil")
	}
	if msg.Request.Method != PushFile {
		t.Errorf("method = %q, want %q", msg.Request.Method, PushFile)
	}
	if msg.OnResponse == nil {
		t.Error("OnResponse channel is nil")
	}

	var stored, exists = socket.Messages[msg.Request.Id]
	if !exists {
		t.Fatal("message not stored in Messages map")
	}
	if stored != msg {
		t.Error("stored message pointer differs from returned")
	}
}

func TestSendClosed(t *testing.T) {
	var socket SSocket

	var msg = socket.send(PushFile, nil)
	if msg != nil {
		t.Error("send() on closed socket should return nil")
	}
}

func TestReceive(t *testing.T) {
	var socket SSocket
	socket.Open()
	t.Cleanup(socket.Close)

	go func() { <-socket.Channel }()

	var msg = socket.send(PushFile, nil)
	if msg == nil {
		t.Fatal("send returned nil")
	}

	var responseJSON = json.RawMessage(
		[]byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":"OK"}`, msg.Request.Id)),
	)

	var got = make(chan bool, 1)
	go func() { got <- <-msg.OnResponse }()

	socket.receive(responseJSON)

	var success = <-got
	if !success {
		t.Error("OnResponse = false, want true")
	}
	if msg.Response == nil {
		t.Fatal("Response not set after receive")
	}
	if msg.Response.Error != nil {
		t.Errorf("unexpected error on response: %v", msg.Response.Error)
	}
}

func TestReceiveError(t *testing.T) {
	var socket SSocket
	socket.Open()
	t.Cleanup(socket.Close)

	go func() { <-socket.Channel }()

	var msg = socket.send(PushFile, nil)
	if msg == nil {
		t.Fatal("send returned nil")
	}

	var responseJSON = json.RawMessage(
		[]byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"error":"not found"}`, msg.Request.Id)),
	)

	var got = make(chan bool, 1)
	go func() { got <- <-msg.OnResponse }()

	socket.receive(responseJSON)

	var success = <-got
	if success {
		t.Error("OnResponse = true, want false for error response")
	}
	if msg.Response == nil {
		t.Fatal("Response not set after receive")
	}
	if msg.Response.Error == nil {
		t.Error("expected Response.Error to be set")
	}
}

func TestReceiveMissingId(t *testing.T) {
	var socket SSocket
	socket.Open()
	t.Cleanup(socket.Close)

	var responseJSON = json.RawMessage(`{"jsonrpc":"2.0","id":999,"result":"OK"}`)

	socket.receive(responseJSON) // should not panic
}
