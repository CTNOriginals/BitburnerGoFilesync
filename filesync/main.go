package main

import (
	"filesync/debug"
	"filesync/server"
)

const port = 8080

func main() {
	print("\n---- FileSync START ----\n")
	defer print("\n---- Program END ----")

	go debug.DebugCommandListener()
	server.Start(port)
}
