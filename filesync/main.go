package main

import (
	"filesync/debug"
	"filesync/rpc"
)

const port = 8080

func main() {
	print("\n---- FileSync START ----\n")
	defer print("\n---- Program END ----")

	go debug.DebugCommandListener()
	rpc.StartServer(port)
}
