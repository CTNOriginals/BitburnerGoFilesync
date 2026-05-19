package websocket

import (
	"encoding/json"
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

type TResponseCallback func(message SMessage)

type SMessage struct {
	Request    *SRequest
	Response   *SResponse
	OnResponse chan bool
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
