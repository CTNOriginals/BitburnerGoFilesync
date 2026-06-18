package watcher

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	"github.com/bmatcuk/doublestar/v4"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")
	testFileEvents()
}

func testFileEvents() {
	var testpath = filepath.Join(config.Values.Directory, "watcher")
	var newPath = filepath.Join(testpath, "new.ext")
	var modPath = filepath.Join(testpath, "mod.ext")
	var delPath = filepath.Join(testpath, "del.ext")

	// clog.Debugf("paths: \n%s\n%s\n%s\n%s\n", testpath, newPath, modPath, delPath)

	var err = os.MkdirAll(testpath, os.ModePerm)
	defer func() {
		var err = os.RemoveAll(testpath)

		if err != nil {
			clog.Error(err)
		}
	}()

	if err != nil {
		clog.Error(err)
	}

	ctnfile.WriteFile(modPath, []string{"not modefied"})
	ctnfile.WriteFile(delPath, []string{"bout to be gone"})

	// time.Sleep(time.Second * 2)
	Initialize()
	go StartScanner()

	ctnfile.WriteFile(newPath, []string{"brand new"})
	ctnfile.WriteFile(modPath, []string{"has been modified"})

	err = os.Remove(delPath)
	if err != nil {
		clog.Error(err)
	}

	time.Sleep(time.Second * 3)
}

func testFilter() {
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

	entry.Recursive(func(entry *SEntry) { entry.UpdateChildList() }, true)
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
