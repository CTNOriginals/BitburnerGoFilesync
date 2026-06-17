package watcher

import (
	"runtime"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

var DirEntry *SEntry

func Initialize() {
	var dir = config.Values.Directory
	clog.Debugf("Watcher Initialize, dir: %s", dir)

	var entry = newEntry(dir)

	if entry.infoError != nil {
		if !entry.Exists() {
			clog.Fatalf(
				"The configured entry directory does not exist: %s\n%v",
				dir, entry.infoError,
			)
		} else {
			clog.Fatalf(
				"An unknown entry error was returned while initializing the entry directory: %s\n%v",
				dir, entry.infoError,
			)
		}

		runtime.Goexit()
	}

	DirEntry = entry
}

func scanner() {
	// TODO:
}
