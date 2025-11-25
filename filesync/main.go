package main

import (
	"filesync/handler"
)

const port = 8080

func main() {
	print("\n---- FileSync START ----\n")
	defer print("\n---- Program END ----")

	println(handler.GetFileByPath("proto.ts"))

	go debug.DebugCommandListener()
	rpc.StartServer(constants.Port)
}
