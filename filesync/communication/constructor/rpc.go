package constructor

import (
	"filesync/communication/definitions"
	"fmt"
)

// The "jsonrpc" field is excluded from this struct, it will be added in RPC.String().
type RPC struct {
	// Id         int
	Method     definitions.Method
	Parameters string
}

func NewRPC(method definitions.Method, parameters ...string) (rpc *RPC) {
	rpc = &RPC{
		Method: method,
	}

	def := definitions.RPCDefinitions[method]

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
		"{\n\"method\": \"%s\",\n \"params\": %s\n}",
		this.Method,
		this.Parameters,
	)
}

func (this RPC) JSON(id int) string {
	return fmt.Sprintf(
		"{\n \"jsonrpc\": \"2.0\",\n \"id\": %d,\n \"method\": \"%s\",\n \"params\": %s\n}",
		id,
		this.Method,
		this.Parameters,
	)
}
