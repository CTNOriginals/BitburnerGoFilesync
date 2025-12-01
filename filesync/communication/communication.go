// Communication stands for the communication between the server and the client.
// The name may be a bit confusing but i was unable to think of a better one for this case.
package communication

import (
	"filesync/communication/constructor"
	"filesync/communication/definitions"
	"filesync/utils"
	"fmt"

	"github.com/gorilla/websocket"
)

// Keeps track of all the outgoing messages by linking them with their id.
// If a response is received, the original request can be found in here as it should share an id.
// Once the response is linked, the response body will be parced and added to the initial message object
var MessageLog constructor.MMessageLog = constructor.MMessageLog{}

var currentId = 0

func GetId() int {
	currentId += 1
	return currentId - 1
}

func SendRequest(method definitions.Method, callback constructor.OnResponseCallback, parameters ...string) *constructor.Message {
	if ActiveConnection == nil {
		fmt.Printf("ActiveConnection is nil\nUnable to send message (%s)\n", method)
		return nil
	}

	var id = GetId()
	var rpc = constructor.NewRPC(method, parameters...)
	var msg = constructor.NewMessage(rpc, callback)
	var json = rpc.JSON(id)

	var err = ActiveConnection.WriteMessage(websocket.TextMessage, []byte(json))

	if err != nil {
		panic(err)
	}

	MessageLog[id] = msg

	return msg
}

func OnResponse(body []byte) {
	var json = utils.JSONToMap(body)
	var id = int(json["id"].(float64))

	if _, exists := MessageLog[id]; !exists {
		fmt.Printf("communication.OnResponse: Unknown id (%d) received.\nDiscarding response:\n%v\n", id, json)
		return
	}

	var message = MessageLog[id]

	if _, exists := json["error"]; exists {
		var val = json["error"]
		fmt.Printf("communication.OnResponse: Message (%d) Received an error response:\n%v\n", id, val)
		message.Response = val
		return
	}

	message.Response = json["result"]

	fmt.Printf("OnResponse %d: %v\n", id, *message)

	message.OnResponse(message)
}
