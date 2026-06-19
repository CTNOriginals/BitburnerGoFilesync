package watcher

import (
	"os"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

var rootDir string
var fileStateMap MFileState

func Initialize() {
	rootDir = config.Values.Directory
	fileStateMap = make(MFileState)
	clog.Debugf("Watcher Initialize, dir: %s", rootDir)

	generatePatternPaths()
	fileStateMap.GetNewEntries(rootDir, func(path string) {
		var err = fileStateMap.Push(path)
		if err != nil {
			clog.Errorf("Encountered an unexpected error while getting file info %s:\n%v", path, err)
		}
	})
}

func StartScanner() {
	<-*websocket.Client.OnReadySub()
	for {
		fileStateMap.GetNewEntries(rootDir, func(path string) {
			var err = fileStateMap.Push(path)

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

func pushFile(path string) {
	websocket.Client.Socket.PushFile(websocket.Params_PushFile{
		Filename: utils.ToBitburnerPath(path),
		Content:  string(utils.GetFileContentByPath(path)),
		Server:   "home",
	}, nil)
}

func onFileCreate(path string) {
	clog.Debugf("On File Create: %s", path)
	pushFile(path)
}
func onFileModify(path string) {
	clog.Debugf("On File Modify: %s", path)
	pushFile(path)
}
func onFileDelete(path string) {
	clog.Debugf("On File Delete: %s", path)
	websocket.Client.Socket.DeleteFile(websocket.Params_DeleteFile{
		Filename: utils.ToBitburnerPath(path),
		Server:   "home",
	}, nil)
}
