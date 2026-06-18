package watcher

import (
	"os"
	"path/filepath"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

var rootDir string
var fileStateMap map[string]time.Time

func Initialize() {
	rootDir = config.Values.Directory
	fileStateMap = make(map[string]time.Time)
	clog.Debugf("Watcher Initialize, dir: %s", rootDir)

	generatePatternPaths()
	getNewEntries(rootDir, func(path string) {
		registerFilePath(path)
	})
}

func registerFilePath(path string) error {
	var info, err = os.Stat(path)

	if err != nil {
		return err
	}

	fileStateMap[path] = info.ModTime()
	return nil
}

func StartScanner() {
	for {
		getNewEntries(rootDir, func(path string) {
			var err = registerFilePath(path)

			if err != nil {
				clog.Errorf("Encountered an unexpected error while getting file info %s:\n%v", path, err)
				return
			}

			onFileCreate(path)
		})

		scan()

		if config.Values.FileScanInterval > 0 {
			time.Sleep(time.Millisecond * time.Duration(config.Values.FileScanInterval))
		}
	}
}

func getNewEntries(dir string, fn func(path string)) {
	var entries, err = os.ReadDir(dir)
	var subDirs = make([]string, 0)

	for _, entry := range entries {
		var path = filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			subDirs = append(subDirs, path)
			continue
		}

		var _, exists = fileStateMap[path]

		if exists || !filePatternFilter(path) {
			continue
		}

		fn(path)
	}

	if err != nil {
		clog.Errorf("Encountered and error while reading directory %s:\n%v", dir, err)
		clog.Debugf("Entries returned before error: %v", entries)
	}

	for _, subDir := range subDirs {
		getNewEntries(subDir, fn)
	}
}

func scan() {
	for path, modTime := range fileStateMap {
		var info, err = os.Stat(path)

		if err != nil {
			if !os.IsNotExist(err) {
				clog.Errorf("Unknown error on path %s:\n%v", path, err)
			}

			onFileDelete(path)
			delete(fileStateMap, path)
			continue
		}

		if modTime.Equal(info.ModTime()) {
			continue
		}

		onFileModify(path)
		fileStateMap[path] = info.ModTime()
	}
}

func onFileCreate(path string) {
	clog.Debugf("TODO: Create file: %s", path)
}
func onFileModify(path string) {
	clog.Debugf("TODO: Push file: %s", path)
}
func onFileDelete(path string) {
	clog.Debugf("TODO: Delete file: %s", path)
}
