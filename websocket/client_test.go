package websocket

import (
	"testing"

	wsgorilla "github.com/gorilla/websocket"
)

func TestClientActive(t *testing.T) {
	var client SClient

	if client.Active() {
		t.Error("Active() = true, want false for nil Connection")
	}

	client.Connection = &wsgorilla.Conn{}
	if !client.Active() {
		t.Error("Active() = false, want true for non-nil Connection")
	}
}
