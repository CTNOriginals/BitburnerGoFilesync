package watcher

import (
	"runtime"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "watcher",
})

var Entry *SEntry

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

	generatePatternPaths()
	entry.SetPathFilter(filePatternFilter)
	entry.UpdateChildList()

	Entry = entry

	clog.Debug(entry.StringRecursive())
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
	Entry.Recursive(func(entry *SEntry) {
		if !entry.IsDirectory() {
			return
		}

		for _, child := range entry.children {
			if child.IsDirectory() || !child.IsModified() {
				continue
			}

			if child.Exists() {
				onFileModify(child)
			} else {
				onFileDelete(child)
			}
		}

		var newEntries = entry.UpdateChildList()

		for _, child := range newEntries {
			onFileModify(child)
		}
	}, true)

	// time.Sleep(time.Second)

	// Update and clean every entry
	// Entry.Recursive((*SEntry).Update, true)
	Entry.Recursive(func(entry *SEntry) {
		// clog.Debugf("Updating %s", entry.info.Name())
		entry.Update()

		entry.CleanChildListFunc(func(child *SEntry) bool {
			child.Update()
			return child.Exists()
		})
	}, true)

	// clog.Message("-- end --")
}

// NOTE: Same as modify
// func onFileCreate(entry *SEntry) {}
func onFileModify(entry *SEntry) {
	clog.Debugf("TODO: Push file: %s", entry.GetPath())
}
func onFileDelete(entry *SEntry) {
	clog.Debugf("TODO: Delete file: %s", entry.GetPath())
}
