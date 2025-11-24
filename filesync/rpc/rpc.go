package rpc

import (
	"fmt"
)

var id = 0

func GetId() int {
	id += 1
	return id - 1
}

type RPC struct {
	Jsonrpc string
	Id      int
	Method  string
	Params  Param
}

func NewRPC(method Method) RPC {
	return RPC{
		Jsonrpc: "2.0",
		Id:      GetId(),
		Method:  string(method),
	}
}

func (this RPC) String() string {
	return fmt.Sprintf(
		"{\n \"jsonrpc\": \"%s\",\n \"id\": %d,\n \"method\": \"%s\",\n \"params\": {%s }\n}",
		this.Jsonrpc,
		this.Id,
		this.Method,
		this.Params.String(),
	)
}
