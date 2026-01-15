package watcher

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

var FileStateMap MFileState = MFileState{}

// Initialize all relevant files and store them
// in a data object along with their current states.
func Initialize() {
	// Register the existing files without calling the OnCreate event
	// to prevent them from being sent over the websocket
	for _, file := range getUnregisteredFiles(config.Values.Directory) {
		FileStateMap[file.Path] = file
	}
}

func FileScanner() {
	fmt.Printf("Scanning files in: %s\n", config.Values.Directory)

	for {
		scanFiles()

		if config.Values.FileScanInterval > 0 {
			time.Sleep(time.Millisecond * time.Duration(config.Values.FileScanInterval))
		}
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

	newFiles := getUnregisteredFiles(config.Values.Directory)

	for _, file := range newFiles {
		FileEventHandlerMap.Handle(file, OnFileCreate)
	}
}

// Checks the dir recursivly for any files
// that are not present in FileStates and returns them.
func getUnregisteredFiles(dir string) (newFiles []*FileInfo) {
	utils.ForEachFileInDirRecursive(dir, func(file os.FileInfo, dir string) {
		if len(config.Values.FilePatterns.Include) > 0 {
			//TODO functionality for wildcard matching (*.js, *.d.ts)
			var split = strings.Split(file.Name(), ".")
			var ext = split[len(split)-1]

			if !slices.Contains(config.Values.FilePatterns.Include, ext) {
				return
			}
		}

		path := fmt.Sprintf("%s/%s", dir, file.Name())
		_, exists := FileStateMap[path]

		if exists {
			return
		}

		newFiles = append(newFiles, &FileInfo{Path: path, Info: file})
	})

	return newFiles
}
