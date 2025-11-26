package rpcHandler

import (
	"filesync/rpc"
	"filesync/rpc/definitions"
	"filesync/utils"

	"github.com/gorilla/websocket"
)

type Message struct {
	Definition definitions.Definition
	Request    rpc.RPC
	Responce   utils.JSONMap
}

func NewMessage(msg rpc.RPC) Message {
	return Message{
		Definition: definitions.RPCDefinitions.GetByMethod(msg.Method),
		Request:    msg,
	}
}

type MMessageLog map[int]Message

var (
	currentId  = 0
	MessageLog MMessageLog
)

func GetId() int {
	currentId += 1
	return currentId - 1
}

func SendRequest(msg rpc.RPC) {
	id := GetId()
	MessageLog[id] = NewMessage(msg)

	rpc.ActiveConnection.WriteMessage(websocket.TextMessage, []byte(msg.JSON(id)))
}

func OnResponse(msg []byte) {}
