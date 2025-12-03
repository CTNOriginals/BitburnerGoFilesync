package watcher

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/communication"
	"github.com/CTNOriginals/BitburnerGoFilesync/communication/constructor"
	"github.com/CTNOriginals/BitburnerGoFilesync/communication/definitions"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

var FileEventHandlerMap = MFileEventHandler{
	OnFileCreate: func(file *FileInfo) {
		fmt.Printf("OnFileCreate: %s\n", file.Path)
		PushFile(file)

		FileStateMap[file.Path] = file
	},
	OnFileModify: func(file *FileInfo) {
		fmt.Printf("OnFileModify: %s\n", file.Path)
		PushFile(file)

		file.Info = file.GetInfo()
	},
	OnFileDelete: func(file *FileInfo) {
		fmt.Printf("OnFileDelete: %s\n", file.Path)

		var relPath = file.RelativePath()

		communication.SendRequest(definitions.DeleteFile, func(message *constructor.Message) {
			if message.IsError {
				fmt.Printf("\nwatcher.OnFileDelete: Unable to delete file: %s\nResponse: %v\n\n", relPath, message.Response)
			}
		}, relPath, "home")

		delete(FileStateMap, file.Path)
	},
}

func PushFile(file *FileInfo) {
	var relPath = file.RelativePath()
	var content = utils.SanitizeFileContent(utils.GetFileContentByPath(relPath))

	communication.SendRequest(definitions.PushFile, func(message *constructor.Message) {
		if message.IsError {
			fmt.Printf("\nwatcher.PushFile: Unable to push file content: %s\nResponse: %v\n\n", relPath, message.Response)
		}
	}, relPath, string(content), "home")
}
