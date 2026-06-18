package watcher

import (
	"os"
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
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

	// time.Sleep(time.Second * 2)
}
