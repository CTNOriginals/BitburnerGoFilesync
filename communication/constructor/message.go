package constructor

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/communication/definitions"

	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

type OnResponseCallback func(message *Message)

type Message struct {
	Definition *definitions.Definition
	Request    *RPC
	Response   any
	OnResponse OnResponseCallback
	IsError    bool
}

func NewMessage(rpc *RPC, callback OnResponseCallback) *Message {
	return &Message{
		Definition: definitions.RPCDefinitions[rpc.Method],
		Request:    rpc,
		OnResponse: callback,
	}
}

func (this Message) String() string {
	return ctnstruct.ToString(this, "OnResponse")
}

type MMessageLog map[int]*Message
