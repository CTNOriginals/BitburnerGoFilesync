package main

import "filesync/server"

const port = 8080

func main() {
	server.Start(port)
}
