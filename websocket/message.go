package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket/rpcschema"

	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

type OnResponseCallback func(message *Message)

type Message struct {
	Schema     *rpcschema.Schema
	Request    *RPC
	Response   any
	OnResponse OnResponseCallback
	IsError    bool
}

func NewMessage(rpc *RPC, callback OnResponseCallback) *Message {
	return &Message{
		Schema:     rpcschema.SchemaMap[rpc.Method],
		Request:    rpc,
		OnResponse: callback,
	}
}

func (this Message) String() string {
	return ctnstruct.ToString(this, "OnResponse")
}

type MMessageLog map[int]*Message
