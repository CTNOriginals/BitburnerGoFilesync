package rpcHandler

import (
	"filesync/rpc"
	"filesync/rpc/definitions"
)

// type TParams1 func
// type TParams2 func(filename string, server string) rpc.RPC

type params1 struct {
	server string
}

type TRequestGenerator func(string) rpc.RPC

var RequestGenerators map[definitions.Method]TRequestGenerator = map[definitions.Method]TRequestGenerator{
	// definitions.GetFile:

	definitions.GetFile: func(server string) rpc.RPC {
		return rpc.NewRPC(definitions.GetFile)
	},
	// definitions.GetAllFiles: func(server string, filename string) rpc.RPC {
	// 	return rpc.NewRPC(definitions.GetAllFiles)
	// },
}

// Parameters: ParameterFields{},
// Parameters: ParameterFields{"server"},
// Parameters: ParameterFields{"filename", "server"},
// Parameters: ParameterFields{"filename", "content", "server"},
