package watcher

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

type EFileEvent int

const (
	OnFileCreate EFileEvent = iota
	OnFileModify
	// TODO: OnFileMove
	// TODO: OnFileRename
	OnFileDelete
)

type MEventHandler map[EFileEvent]func(path string)

func (this MEventHandler) Emit(event EFileEvent, path string) {
	this[event](path)
}

var fileEventHandler = MEventHandler{
	OnFileCreate: func(path string) {
		clog.Infof("On File Create: %s", path)
		pushFile(path)
	},
	OnFileModify: func(path string) {
		clog.Infof("On File Modify: %s", path)
		pushFile(path)
	},
	OnFileDelete: func(path string) {
		clog.Infof("On File Delete: %s", path)
		websocket.Client.Socket.DeleteFile(websocket.Params_DeleteFile{
			Filename: utils.ToBitburnerPath(path),
			Server:   "home",
		}, nil)
	},
}

func pushFile(path string) {
	websocket.Client.Socket.PushFile(websocket.Params_PushFile{
		Filename: utils.ToBitburnerPath(path),
		Content:  string(utils.GetFileContentByPath(path)),
		Server:   "home",
	}, nil)
}
