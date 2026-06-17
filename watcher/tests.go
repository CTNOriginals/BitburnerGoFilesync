package watcher

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/bmatcuk/doublestar/v4"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	var baseDir = config.Values.Directory
	var entry = newEntry(baseDir)
	// var entry = newEntry(constants.WorkindDirectory)

	if entry.infoError != nil {
		clog.Error(entry.infoError)
	}

	var inc = filepath.Join(baseDir, "**/*.ext")
	var exc = filepath.Join(baseDir, "**/*.x.ext")
	clog.Debugf("inc: %s", inc)
	clog.Debugf("exc: %s", exc)

	var incList, _ = doublestar.FilepathGlob(inc)
	var excList, _ = doublestar.FilepathGlob(exc)
	clog.Debugf("inc list:\n%s", strings.Join(incList, "\n"))
	clog.Debugf("exc list:\n%s", strings.Join(excList, "\n"))

	entry.Recursive((*SEntry).UpdateChildList, true)
	clog.Infof("%s:\n%s", entry.GetPath(), entry.StringRecursive())

	var filter = func(path string, info os.FileInfo) bool {
		if info.IsDir() {
			return true
		}

		clog.Debugf("path: %s", path)

		if slices.Contains(excList, path) {
			return false
		}

		// Also check if include is not nil, and if it is, just return true.

		if !slices.Contains(incList, path) {
			return false
		}

		return true
	}

	entry.Recursive(func(child *SEntry) {
		child.SetPathFilter(filter)
	}, true)

	entry.Recursive((*SEntry).ApplyPathFilter, true)
	clog.Infof("%s:\n%s", entry.GetPath(), entry.StringRecursive())
}
