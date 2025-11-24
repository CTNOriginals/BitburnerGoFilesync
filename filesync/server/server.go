package server

import (
	"filesync/rpc"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(port int) {
	print("\n---- Starting Server ----\n")

	http.HandleFunc("/", wsHandler)
	fmt.Printf("WebSocket server started on: %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	print("Home connected!\n")

	testMsg := rpc.NewRPC(rpc.GetFileNames)
	testMsg.Params.Add("server", "home")

	fmt.Printf("Sending test message: %s\n", testMsg.String())

	err = conn.WriteMessage(websocket.TextMessage, []byte(testMsg.String()))
	if err != nil {
		println(err)
		return
	}

	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		// // Echo the message back to the client
		// if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// 	fmt.Println("Error writing message:", err)
		// 	break
		// }
	}
}
