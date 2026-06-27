package websocket

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	var req = NewRequest(42, PushFile, "testparam")

	if req.Jsonrpc != "2.0" {
		t.Errorf("Jsonrpc = %q, want %q", req.Jsonrpc, "2.0")
	}
	if req.Id != 42 {
		t.Errorf("Id = %d, want %d", req.Id, 42)
	}
	if req.Method != PushFile {
		t.Errorf("Method = %q, want %q", req.Method, PushFile)
	}
	if req.Params.(string) != "testparam" {
		t.Errorf("Params = %v, want %v", req.Params, "testparam")
	}
}

func TestSRequestString(t *testing.T) {
	var req = NewRequest(1, PushFile, "test.js")
	var want = "pushFile(1): test.js"
	if s := req.String(); s != want {
		t.Errorf("String() = %q, want %q", s, want)
	}
}

func TestSResponseStringWithError(t *testing.T) {
	var resp = SResponse{Id: 1, Error: "not found"}
	var s = resp.String()
	if !strings.Contains(s, "ERROR") {
		t.Errorf("String() = %q, want it to contain 'ERROR'", s)
	}
}

func TestSResponseStringWithResult(t *testing.T) {
	var resp = SResponse{Id: 1, Result: json.RawMessage(`"OK"`)}
	var s = resp.String()
	if len(s) == 0 {
		t.Fatal("String() returned empty")
	}
}

func TestSMessageStringRequestOnly(t *testing.T) {
	var msg = SMessage{Request: NewRequest(1, PushFile, nil)}
	var s = msg.String()
	if !strings.HasPrefix(s, "Request: ") {
		t.Errorf("String() = %q, want 'Request: ' prefix", s)
	}
}

func TestSMessageStringWithResponse(t *testing.T) {
	var msg = SMessage{
		Request:  NewRequest(1, PushFile, nil),
		Response: &SResponse{Id: 1, Result: json.RawMessage(`"OK"`)},
	}
	var s = msg.String()
	if !strings.Contains(s, "Response:") {
		t.Errorf("String() = %q, want it to contain 'Response:'", s)
	}
}
