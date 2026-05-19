package websocket

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SRequest struct {
	Jsonrpc string  `json:"jsonrpc"`
	Id      int     `json:"id"`
	Method  TMethod `json:"method"`
	Params  any     `json:"params,omitempty"`
}
type SResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   any             `json:"error,omitempty"`
}

func NewRequest(id int, method TMethod, params any) *SRequest {
	return &SRequest{
		Jsonrpc: "2.0",
		Id:      id,
		Method:  method,
		Params:  params,
	}
}

func (this SRequest) String() string {
	return fmt.Sprintf("%s(%d): %v", this.Method, this.Id, this.Params)
}
func (this SResponse) String() string {
	if this.Error != nil {
		return fmt.Sprintf("%d: ERROR {\n%v\n}", this.Id, this.Error)
	}

	return fmt.Sprintf("%d: %v", this.Id, this.Result.UnmarshalJSON(this.Result))
}

type TResponseCallback func(message SMessage)

type SMessage struct {
	Request    *SRequest
	Response   *SResponse
	OnResponse chan bool
}

func (this SMessage) String() string {
	var str strings.Builder

	str.WriteString("Request: ")
	str.WriteString(this.Request.String())

	if this.Response == nil {
		return str.String()
	}

	str.WriteString("Response: ")
	str.WriteString(this.Response.String())

	return str.String()
}

// func (this *SMessage) Receive(body json.RawMessage) error {
// 	var response *SResponse = nil
// 	var err = json.Unmarshal(body, response)
//
// 	if err != nil {
// 		return fmt.Errorf("Message.Receive error while trying to unmarshal body: %v", err)
// 	}
//
// 	if response.Error != nil {
// 		this.OnResponse <- false
// 		return fmt.Errorf("%v", response.Error)
// 	}
//
// 	this.OnResponse <- true
//
// 	return nil
// }
