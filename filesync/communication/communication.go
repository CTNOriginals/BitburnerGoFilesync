// Communication stands for the communication between the server and the client.
// The name may be a bit confusing but i was unable to think of a better one for this case.
package communication

import (
	"filesync/communication/constructor"
	"filesync/communication/definitions"
	"fmt"

	"github.com/gorilla/websocket"
)

var currentId = 0

// Keeps track of all the outgoing messages by linking them with their id.
// If a response is received, the original request can be found in here as it should share an id.
// Once the response is linked, the response body will be parced and added to the initial message object
var MessageLog constructor.MMessageLog = constructor.MMessageLog{}

func GetId() int {
	currentId += 1
	return currentId - 1
}

func SendRequest(method definitions.Method, parameters ...string) *constructor.Message {
	if ActiveConnection == nil {
		fmt.Printf("ActiveConnection is nil\nUnable to send message (%s)\n", method)
		return nil
	}

	var id = GetId()
	var rpc = constructor.NewRPC(method, parameters...)
	var msg = constructor.NewMessage(rpc)
	var json = rpc.JSON(id)

	var err = ActiveConnection.WriteMessage(websocket.TextMessage, []byte(json))

	if err != nil {
		println(err)
		return nil
	}

	MessageLog[id] = *constructor.NewMessage(rpc)

	return msg
}

func OnResponse(msg []byte) {}
