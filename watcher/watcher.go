package watcher

import (
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

func Initialize() {
	var dir = config.Values.Directory
	clog.Debugf("Watcher Initialize, dir: %s", dir)
}

func StartScanner() {
	for {
		scan()

		if config.Values.FileScanInterval > 0 {
			time.Sleep(time.Millisecond * time.Duration(config.Values.FileScanInterval))
		}
	}
}

func scan() {

}

func onFileModify(file *SFileData) {
	clog.Debugf("TODO: Push file: %s", file.Path)
}
func onFileDelete(file *SFileData) {
	clog.Debugf("TODO: Delete file: %s", file.Path)
}
