package main

import (
	"filesync/constants"
	"filesync/debug"
	"filesync/rpc"
)

func main() {
	print("\n\n---- FileSync START ----\n")
	defer print("---- Program END ----")

	// fileContent := handler.GetFileByPath("fileTest.ts")
	// sanitized := handler.SanitizeFileContent(fileContent)
	// println(string(sanitized))

	go debug.DebugCommandListener()
	rpc.StartServer(constants.Port)
}
