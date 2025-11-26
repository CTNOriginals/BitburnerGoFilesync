package watcher

import (
	"filesync/constants"
	"filesync/handler"
	"fmt"
	"os"
	"time"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/utils/file"
)

var FileStateMap MFileState = MFileState{}

var FileEventHandlerMap = MFileEventHandler{
	OnFileCreate: func(file *FileInfo) {
		fmt.Printf("%s: OnFileCreate\n", file.Path)
		// TODO Push the file over the websocket
		FileStateMap[file.Path] = file
	},
	OnFileModify: func(file *FileInfo) {
		fmt.Printf("%s: OnFileModify\n", file.Path)
		// TODO Handle websocket
		file.Info = file.GetInfo()
	},
	OnFileDelete: func(file *FileInfo) {
		fmt.Printf("%s: OnFileDelete\n", file.Path)
		// TODO Handle websocket
		delete(FileStateMap, file.Path)
	},
}

// Initialize all relevant files and store them
// in a data object along with their current states.
func Initialize() {
	// Register the existing files without calling the OnCreate event
	// to prevent them from being sent over the websocket
	for _, file := range getUnregisteredFiles(constants.BitburnerRelativePath) {
		FileStateMap[file.Path] = file
	}
}

func FileScanner() {
	for {
		scanFiles()
		time.Sleep(time.Millisecond * 100)
	}
}

// Scans all file states in FileStates.
// Once a change is detected on one of the files,
// the event will be added to the state.
func scanFiles() {
	for path, state := range FileStateMap {
		if !ctnfile.FileExists(path) {
			FileEventHandlerMap.Handle(state, OnFileDelete)
			continue
		}

		if state.GetInfo().ModTime().Compare(state.Info.ModTime()) == 1 {
			FileEventHandlerMap.Handle(state, OnFileModify)
		}
	}

	newFiles := getUnregisteredFiles(constants.BitburnerRelativePath)

	for _, file := range newFiles {
		FileEventHandlerMap.Handle(file, OnFileCreate)
	}
}

// Checks the dir recursivly for any files
// that are not present in FileStates and returns them.
func getUnregisteredFiles(dir string) (newFiles []*FileInfo) {
	handler.ForEachFileInDirRecursive(dir, func(file os.FileInfo, dir string) {
		path := fmt.Sprintf("%s/%s", dir, file.Name())
		_, exists := FileStateMap[path]

		if exists {
			return
		}

		newFiles = append(newFiles, &FileInfo{Path: path, Info: file})
	})

	return newFiles
}
