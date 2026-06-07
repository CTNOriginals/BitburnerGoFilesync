package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/bmatcuk/doublestar/v4"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

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
	clog.Infof("Scanning files in: %s\n", config.Values.Directory)

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
		var fullPath = utils.GetAbsolutePath(path)

		if !ctnfile.FileExists(fullPath) {
			FileEventHandlerMap.Handle(state, OnFileDelete)
			continue
		}

		if state.GetInfo().ModTime().Compare(state.Info.ModTime()) == 1 {
			FileEventHandlerMap.Handle(state, OnFileModify)
		}
	}

	newFiles := getUnregisteredFiles(config.Values.Directory)

	for _, file := range newFiles {
		clog.Debugf("new file: %s\n", file)
		FileEventHandlerMap.Handle(file, OnFileCreate)
	}
}

// Checks the dir recursivly for any files
// that are not present in FileStates and returns them.
func getUnregisteredFiles(dir string) (newFiles []*FileInfo) {
	utils.ForEachFileInDirRecursive(dir, func(file os.FileInfo, dir string) {
		var reldir = strings.Replace(dir, config.Values.Directory, "", 1)
		// Relative path to bitburners root dir
		var path string

		if reldir == "" {
			path = file.Name()
		} else {
			path = fmt.Sprintf("%s/%s", reldir, file.Name())

			// Remove the first '/' if it is present.
			// This ensures a little more consistency
			// revative to the paths that are at the root level.
			if path[0] == '/' {
				path = path[1:]
			}
		}

		if !shouldIncludeFile(path) {
			return
		}

		_, exists := FileStateMap[path]

		if exists {
			return
		}

		newFiles = append(newFiles, &FileInfo{Path: path, Info: file})
	})

	return newFiles
}

// Check if the file path should be included
// according to the config values Include and Exclude patternd
func shouldIncludeFile(path string) bool {
	for _, pattern := range config.Values.FilePatterns.Exclude {
		if patternMatch(pattern, path) {
			return false
		}
	}

	for _, pattern := range config.Values.FilePatterns.Include {
		if patternMatch(pattern, path) {
			return true
		}
	}

	return len(config.Values.FilePatterns.Include) == 0
}

//	func patternMatch(pattern string, path string) bool {
//		var match, err = doublestar.PathMatch(pattern, path)
//
//		if err != nil {
//			clog.Errorf(
//				"Pattern match error: %v (%s > %s = %t)\n",
//				err,
//				pattern,
//				path,
//				match,
//			)
//			return false
//		}
//
//		return match
//	}
func patternMatch(pattern string, path string) bool {
	var patternPath = fmt.Sprintf("%s%s%s", config.Values.Directory, string(filepath.Separator), pattern)

	var match, err = doublestar.FilepathGlob(patternPath)
	// var match, err = filepath.Glob(patternPath)

	if err != nil {
		clog.Errorf(
			"Pattern match error: %v (%s > %s = %t)\n",
			err,
			patternPath,
			path,
			slices.Contains(match, path),
		)
		return false
	}

	if len(match) > 0 {
		clog.Debugf("Pattern %s returns: \n%s", patternPath, strings.Join(match, "\n"))
	}

	return slices.Contains(match, utils.GetAbsolutePath(path))
}
