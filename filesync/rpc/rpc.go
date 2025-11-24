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
	Id     int
	Method string
	Params Param
}

func NewRPC(method Method) RPC {
	return RPC{
		Id:     GetId(),
		Method: string(method),
	}
}

func (this RPC) String() string {
	return fmt.Sprintf(
		"{\n \"jsonrpc\": \"2.0\",\n \"id\": %d,\n \"method\": \"%s\",\n \"params\": {%s }\n}",
		this.Id,
		this.Method,
		this.Params.String(),
	)
}
