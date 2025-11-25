package main

import (
	"filesync/constants"
	"filesync/debug"
	"filesync/rpc"
)

func main() {
	print("\n\n---- FileSync START ----\n")
	defer print("---- Program END ----")

	test()

	go debug.DebugCommandListener()
	rpc.StartServer(constants.Port)
}

func test() {
	// // fileContent := handler.GetFileByPath("fileTest.ts")
	// // sanitized := handler.SanitizeFileContent(fileContent)
	// // println(string(sanitized))

	// def := definitions.RPCDefinitions.GetByMethod(definitions.GetFile)
	// defaultParameters, _ := debug.DefaultParameters[def.Method]

	// jsonRpc := rpc.NewRPC(def.Method, defaultParameters...).String()
	// jsonMap := handler.JSONToMap([]byte(jsonRpc))

	// println(jsonRpc)

	// for key, val := range *jsonMap {
	// 	fmt.Printf("%s: (%T) %v\n", key, val, val)
	// }
}
