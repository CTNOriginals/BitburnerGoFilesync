package constructor

import (
	"filesync/communication/definitions"
	"filesync/utils"

	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
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

func (this Message) String() string {
	return ctnstruct.ToString(this)
}

func (this Message) Result() any {
	return this.Responce["result"]
}

// func (this Message) GetResponseField(field string) any {
// 	for _, v := range v {

// 	}
// }

type MMessageLog map[int]Message
