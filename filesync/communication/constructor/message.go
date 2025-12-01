package constructor

import (
	"filesync/communication/definitions"

	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

type Message struct {
	Definition *definitions.Definition
	Request    *RPC
	Response   any
}

func NewMessage(rpc *RPC) *Message {
	return &Message{
		Definition: definitions.RPCDefinitions[rpc.Method],
		Request:    rpc,
	}
}

func (this Message) String() string {
	return ctnstruct.ToString(this)
}

type MMessageLog map[int]*Message
