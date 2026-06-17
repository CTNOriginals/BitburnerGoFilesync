package watcher

import (
	"os"
	"runtime"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

var Entry *SEntry
var Directories *SEntry

func Initialize() {
	var dir = config.Values.Directory
	clog.Debugf("Watcher Initialize, dir: %s", dir)

	var entry = newEntry(dir)
	var dirEntries = newEntry(dir)

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

	generatePatternPaths()
	entry.SetPathFilter(filePatternFilter)
	dirEntries.SetPathFilter(func(_ string, info os.FileInfo) bool {
		return info.IsDir()
	})

	entry.UpdateChildList()
	dirEntries.UpdateChildList()

	Entry = entry
	Directories = dirEntries

	clog.Debug(entry.StringRecursive())
	clog.Debug(dirEntries.StringRecursive())
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
	var modified = make([]*SEntry, 0)

	Entry.Recursive(func(entry *SEntry) {
		if (entry.IsDirectory() && entry.IsModified()) == false {
			return
		}

		entry.ForEachChild(func(child *SEntry) {
			// clog.Debugf("%s > %s", entry.String(), child.String())
			if child.IsModified() {
				modified = append(modified, child)
			}
		})
	}, true)

	clog.Debug("----\n ")

	if len(modified) == 0 {
		return
	}

	for _, file := range modified {
		clog.Debugf("Modified: %s", file.GetPath())
	}
}
