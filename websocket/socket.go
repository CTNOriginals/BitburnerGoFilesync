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

	// log.Printf("Socket.send sending message: %v\n", *message)

	this.Messages[id] = message
	this.Channel <- message

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

func AwaitResponse[T any](message *SMessage) *T {
	var success = <-message.OnResponse

	if success == false {
		log.Printf("Socket.AwaitResponse received response error: %v\n", message.Response)
		return nil
	}

	var result T
	var err = json.Unmarshal(message.Response.Result, &result)

	if err != nil {
		log.Printf("Socket.AwaitResponse error while parsing response result: %v\n", err)
		return nil
	}

	return &result
}

func (this *SSocket) PushFile(
	params Params_PushFile,
	callback func(result *Result_PushFile),
) {
	var message = this.send(PushFile, params)
	callback(AwaitResponse[Result_PushFile](message))
}

func (this *SSocket) GetFile(
	params Params_GetFile,
	callback func(result *Result_GetFile),
) {
	var message = this.send(GetFile, params)
	callback(AwaitResponse[Result_GetFile](message))
}

func (this *SSocket) GetFileMetadata(
	params Params_GetFileMetadata,
	callback func(result *Result_GetFileMetadata),
) {
	var message = this.send(GetFileMetadata, params)
	callback(AwaitResponse[Result_GetFileMetadata](message))
}

func (this *SSocket) DeleteFile(
	params Params_DeleteFile,
	callback func(result *Result_DeleteFile),
) {
	var message = this.send(DeleteFile, params)
	callback(AwaitResponse[Result_DeleteFile](message))
}

func (this *SSocket) GetFileNames(
	params Params_GetFileNames,
	callback func(result *Result_GetFileNames),
) {
	var message = this.send(GetFileNames, params)
	callback(AwaitResponse[Result_GetFileNames](message))
}

func (this *SSocket) GetAllFiles(
	params Params_GetAllFiles,
	callback func(result *Result_GetAllFiles),
) {
	var message = this.send(GetAllFiles, params)
	callback(AwaitResponse[Result_GetAllFiles](message))
}

func (this *SSocket) GetAllFileMetadata(
	params Params_GetAllFileMetadata,
	callback func(result *Result_GetAllFileMetadata),
) {
	var message = this.send(GetAllFileMetadata, params)
	callback(AwaitResponse[Result_GetAllFileMetadata](message))
}

func (this *SSocket) CalculateRam(
	params Params_CalculateRam,
	callback func(result *Result_CalculateRam),
) {
	var message = this.send(CalculateRam, params)
	callback(AwaitResponse[Result_CalculateRam](message))
}

func (this *SSocket) GetDefinitionFile(
	callback func(result *Result_GetDefinitionFile),
) {
	var message = this.send(GetDefinitionFile, nil)
	callback(AwaitResponse[Result_GetDefinitionFile](message))
}

func (this *SSocket) GetSaveFile(
	callback func(result *Result_GetSaveFile),
) {
	var message = this.send(GetSaveFile, nil)
	callback(AwaitResponse[Result_GetSaveFile](message))
}

func (this *SSocket) GetAllServers(
	callback func(result *Result_GetAllServers),
) {
	var message = this.send(GetAllServers, nil)
	callback(AwaitResponse[Result_GetAllServers](message))
}
