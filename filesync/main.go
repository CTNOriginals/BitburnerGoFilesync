package main

import (
	"bufio"
	"filesync/server"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

const port = 8080

func main() {
	print("\n---- FileSync START ----\n")
	defer print("\n---- Program END ----")

	go commandListener()
	server.Start(port)
}

func commandListener() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())

		if cmd == "" {
			continue
		}

		if server.ActiveConnection == nil {
			fmt.Println("No active connection")
			continue
		}

		// Send command as message
		err := server.ActiveConnection.WriteMessage(websocket.TextMessage, []byte(cmd))
		if err != nil {
			fmt.Println("Error sending:", err)
		}
	}
}
