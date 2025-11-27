package test

import (
	"filesync/communication/constructor"
	"filesync/debug"
)

func DoTest() {
	for method, params := range debug.DefaultParameters {
		var rpc = constructor.NewRPC(method, params...)
		println(rpc.String())
	}
}
