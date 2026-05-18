package watcher

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket/definitions"
)

var FileEventHandlerMap = MFileEventHandler{
	OnFileCreate: func(file *FileInfo) {
		log.Printf("OnFileCreate: %s\n", file.Path)
		PushFile(file)

		FileStateMap[file.Path] = file
	},
	OnFileModify: func(file *FileInfo) {
		log.Printf("OnFileModify: %s\n", file.Path)
		PushFile(file)

		file.Info = file.GetInfo()
	},
	OnFileDelete: func(file *FileInfo) {
		log.Printf("OnFileDelete: %s\n", file.Path)

		var relPath = file.RelativePath()

		websocket.SendRequest(definitions.DeleteFile, func(message *websocket.Message) {
			if message.IsError {
				log.Printf("\nwatcher.OnFileDelete: Unable to delete file: %s\nResponse: %v\n\n", relPath, message.Response)
			}
		}, relPath, "home")

		delete(FileStateMap, file.Path)
	},
}

func PushFile(file *FileInfo) {
	var relPath = file.RelativePath()
	var content = utils.SanitizeFileContent(utils.GetFileContentByPath(relPath))

	websocket.SendRequest(definitions.PushFile, func(message *websocket.Message) {
		if message.IsError {
			log.Printf("\nwatcher.PushFile: Unable to push file content: %s\nResponse: %v\n\n", relPath, message.Response)
		}
	}, relPath, string(content), "home")
}
