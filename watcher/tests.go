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
	var node = newNode(baseDir)
	// var node = newNode(constants.WorkindDirectory)

	if node.infoError != nil {
		clog.Error(node.infoError)
	}

	var inc = filepath.Join(baseDir, "**/*.ext")
	var exc = filepath.Join(baseDir, "**/*.x.ext")
	clog.Debugf("inc: %s", inc)
	clog.Debugf("exc: %s", exc)

	var incList, _ = doublestar.FilepathGlob(inc)
	var excList, _ = doublestar.FilepathGlob(exc)
	clog.Debugf("inc list:\n%s", strings.Join(incList, "\n"))
	clog.Debugf("exc list:\n%s", strings.Join(excList, "\n"))

	node.Recursive((*SNode).UpdateChildList, true)
	clog.Infof("%s:\n%s", node.GetPath(), node.StringRecursive())

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

	node.Recursive(func(child *SNode) {
		child.SetPathFilter(filter)
	}, true)

	node.Recursive((*SNode).ApplyPathFilter, true)
	clog.Infof("%s:\n%s", node.GetPath(), node.StringRecursive())
}
