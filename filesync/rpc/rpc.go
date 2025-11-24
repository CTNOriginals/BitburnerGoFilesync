package rpc

import (
	"filesync/rpc/definitions"
	"fmt"
)

var id = 0

func GetId() int {
	id += 1
	return id - 1
}

// The "jsonrpc" field is excluded from this struct, it will be added in RPC.String().
type RPC struct {
	Id         int
	Method     string
	Parameters string
}

func NewRPC(method definitions.Method, parameters ...string) (rpc RPC) {
	rpc = RPC{
		Id:     GetId(),
		Method: string(method),
	}

	def := definitions.RPCDefinitions.GetByMethod(method)

	if def.IsError() {
		print(fmt.Errorf("RPC.NewRPC: Invalid parameter (method): %s\n", method))
		return
	}

	if len(def.Parameters) > 0 {
		rpc.Parameters = def.Parameters.Generate(parameters)
	}

	return rpc
}

func (this RPC) String() string {
	return fmt.Sprintf(
		"{\n \"jsonrpc\": \"2.0\",\n \"id\": %d,\n \"method\": \"%s\",\n \"params\": %s\n}",
		this.Id,
		this.Method,
		this.Parameters,
	)
}
