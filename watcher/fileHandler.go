package watcher

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var FileEventHandlerMap = MFileEventHandler{
	OnFileCreate: func(file *FileInfo) {
		clog.Infof("OnFileCreate: %s\n", file.Path)
		PushFile(file)

		FileStateMap[file.Path] = file
	},
	OnFileModify: func(file *FileInfo) {
		clog.Infof("OnFileModify: %s\n", file.Path)
		PushFile(file)

		file.Info = file.GetInfo()
	},
	OnFileDelete: func(file *FileInfo) {
		clog.Infof("OnFileDelete: %s\n", file.Path)

		var relPath = file.RelativePath()

		websocket.Client.Socket.DeleteFile(websocket.Params_DeleteFile{
			Filename: relPath,
			Server:   "home",
		}, nil)

		delete(FileStateMap, file.Path)
	},
}

func PushFile(file *FileInfo) {
	var relPath = file.RelativePath()
	// var content = utils.SanitizeFileContent(utils.GetFileContentByPath(relPath))
	var content = utils.GetFileContentByPath(relPath)

	websocket.Client.Socket.PushFile(websocket.Params_PushFile{
		Filename: relPath,
		Content:  string(content),
		Server:   "home",
	}, nil)
}
