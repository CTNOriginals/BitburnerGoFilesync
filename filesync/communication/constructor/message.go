package constructor

import (
	"filesync/communication/definitions"
	"filesync/utils"
)

type Message struct {
	Definition definitions.Definition
	Request    RPC
	Responce   utils.JSONMap
}

func NewMessage(msg RPC) Message {
	return Message{
		Definition: definitions.RPCDefinitions[msg.Method],
		Request:    msg,
	}
}

type MMessageLog map[int]Message
