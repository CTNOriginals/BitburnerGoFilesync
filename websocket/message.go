package websocket

type SMessage struct {
	Id         int
	Request    any // TODO:
	Response   any // TODO:
	OnResponse func(message SMessage)
}
