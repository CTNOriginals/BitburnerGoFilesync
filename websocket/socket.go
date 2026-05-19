package websocket

import (
	"encoding/json"
	"log"
)

type SSocket struct {
	Channel  chan *SMessage
	Messages map[int]*SMessage

	currentId int
}

func (this *SSocket) Open() {
	this.Channel = make(chan *SMessage)
	this.Messages = map[int]*SMessage{}
}
func (this *SSocket) Close() {
	close(this.Channel)
	this.Channel = nil
}

func (this *SSocket) getId() int {
	var id = this.currentId
	this.currentId += 1
	return id
}

func (this *SSocket) send(method TMethod, params any) *SMessage {
	var id = this.getId()
	var message = &SMessage{
		Request:    NewRequest(id, method, params),
		OnResponse: make(chan bool),
	}

	log.Printf("Socket.send sending message: %v\n", *message)

	this.Channel <- message
	this.Messages[id] = message

	return message
}

func (this *SSocket) Receive(body json.RawMessage) {
	var response SResponse
	var err = json.Unmarshal(body, &response)

	if err != nil {
		log.Printf("Socket.Receive error while trying to unmarshal body: %v\n", err)
		return
	}

	var message, exists = this.Messages[response.Id]

	if !exists {
		log.Printf("Socket.Receive message id does not exist: %d\n", response.Id)
		return
	}

	message.Response = &response

	if response.Error != nil {
		message.OnResponse <- false
		log.Printf("Socket.Receive response contained error: %v\n", response.Error)
		return
	}

	message.OnResponse <- true
}

type getAllFilesParams struct {
	Server string `json:"server"`
}
type getAllFilesResult struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (this *SSocket) GetAllFiles(
	params getAllFilesParams,
	callback func(result []getAllFilesResult),
) {
	var message = this.send(GetAllFiles, params)

	var success = <-message.OnResponse

	if success == false {
		log.Printf("Socket.GetAllFiles reveived response error: %v\n", message.Response)
		return
	}

	var result []getAllFilesResult
	var err = json.Unmarshal(message.Response.Result, &result)

	if err != nil {
		log.Printf("Socket.GetAllFiles error while parsing response result: %v\n", err)
		return
	}

	callback(result)
}
