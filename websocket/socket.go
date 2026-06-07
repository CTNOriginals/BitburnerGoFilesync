package websocket

import (
	"encoding/json"
	"sync"
)

type SSocket struct {
	Channel  chan *SMessage
	Messages map[int]*SMessage

	currentId int
	isOpen    bool

	mutex sync.Mutex
}

func (this *SSocket) Open() {
	if this.isOpen {
		clog.Errorf("Attempted to open an alrady open socket:\n%s", clog.GetStackTrace(0))
		return
	}

	this.Channel = make(chan *SMessage)
	this.Messages = map[int]*SMessage{}
	this.isOpen = true
}
func (this *SSocket) Close() {
	if !this.isOpen {
		clog.Errorf("Attempted to close an alrady closed socket:\n%s", clog.GetStackTrace(0))
		return
	}

	close(this.Channel)
	this.Channel = nil
	this.isOpen = false
}

func (this *SSocket) getId() int {
	var id = this.currentId
	this.currentId += 1
	return id
}

func (this *SSocket) send(method TMethod, params any) *SMessage {
	this.mutex.Lock()

	if !this.isOpen {
		clog.Error("Unable to send message while socket is closed.")
		clog.Debug("TODO: make buffer for messages to send once socket opens.")
		return nil
	}

	var id = this.getId()
	var message = &SMessage{
		Request:    NewRequest(id, method, params),
		OnResponse: make(chan bool),
	}

	clog.Debugf("sending message: %v\n", *message)

	this.Messages[id] = message
	this.Channel <- message

	this.mutex.Unlock()

	return message
}

func (this *SSocket) receive(body json.RawMessage) {
	this.mutex.Lock()

	var response SResponse
	var err = json.Unmarshal(body, &response)

	if err != nil {
		clog.Fatalf("Error while trying to unmarshal body: %v\n", err)
		return
	}

	var message, exists = this.Messages[response.Id]
	clog.Debugf("Received message: %s", string(body))

	if !exists {
		clog.Errorf("Message id does not exist: %d\n", response.Id)
		return
	}

	message.Response = &response

	if response.Error != nil {
		message.OnResponse <- false
		clog.Errorf("Response contained error: %v\n", response.Error)
		return
	}

	message.OnResponse <- true

	this.mutex.Unlock()
}

func AwaitResponse[T any](message *SMessage) *T {
	var success = <-message.OnResponse

	if success == false {
		clog.Errorf("Received response error: %v\n", message.Response)
		return nil
	}

	var result T
	var err = json.Unmarshal(message.Response.Result, &result)

	if err != nil {
		clog.Errorf("Error while parsing response result: %v\n", err)
		return nil
	}

	return &result
}

func call[TResult any](this *SSocket, method TMethod, params any, callback func(*TResult)) {
	var message = this.send(method, params)
	if callback == nil {
		return
	}
	callback(AwaitResponse[TResult](message))
}

func (this *SSocket) PushFile(p Params_PushFile, cb func(*Result_PushFile)) {
	call(this, PushFile, p, cb)
}
func (this *SSocket) GetFile(p Params_GetFile, cb func(*Result_GetFile)) {
	call(this, GetFile, p, cb)
}
func (this *SSocket) GetFileMetadata(p Params_GetFileMetadata, cb func(*Result_GetFileMetadata)) {
	call(this, GetFileMetadata, p, cb)
}
func (this *SSocket) DeleteFile(p Params_DeleteFile, cb func(*Result_DeleteFile)) {
	call(this, DeleteFile, p, cb)
}
func (this *SSocket) GetFileNames(p Params_GetFileNames, cb func(*Result_GetFileNames)) {
	call(this, GetFileNames, p, cb)
}
func (this *SSocket) GetAllFiles(p Params_GetAllFiles, cb func(*Result_GetAllFiles)) {
	call(this, GetAllFiles, p, cb)
}
func (this *SSocket) GetAllFileMetadata(p Params_GetAllFileMetadata, cb func(*Result_GetAllFileMetadata)) {
	call(this, GetAllFileMetadata, p, cb)
}
func (this *SSocket) CalculateRam(p Params_CalculateRam, cb func(*Result_CalculateRam)) {
	call(this, CalculateRam, p, cb)
}
func (this *SSocket) GetDefinitionFile(cb func(*Result_GetDefinitionFile)) {
	call(this, GetDefinitionFile, nil, cb)
}
func (this *SSocket) GetSaveFile(cb func(*Result_GetSaveFile)) {
	call(this, GetSaveFile, nil, cb)
}
func (this *SSocket) GetAllServers(cb func(*Result_GetAllServers)) {
	call(this, GetAllServers, nil, cb)
}
