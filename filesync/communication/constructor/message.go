package constructor

import (
	"filesync/communication/definitions"
	"filesync/utils"
)

type Message struct {
	Definition definitions.Definition
	Request    *RPC
	Responce   utils.JSONMap
}

func NewMessage(rpc *RPC) *Message {
	return &Message{
		Definition: definitions.RPCDefinitions[rpc.Method],
		Request:    rpc,
	}
}

type MMessageLog map[int]Message
