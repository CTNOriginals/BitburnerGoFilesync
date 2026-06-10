package watcher

import (
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
