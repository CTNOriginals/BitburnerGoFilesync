package watcher

import (
	"os"
	"path/filepath"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")
	testFileEvents()
}

func testFileEvents() {
	var testpath = filepath.Join(config.Values.Directory, "watcher")
	var newPath = filepath.Join(testpath, "new.ts")
	var newExPath = filepath.Join(testpath, "new.d.ts")
	var modPath = filepath.Join(testpath, "mod.ts")
	var delPath = filepath.Join(testpath, "del.ts")

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

	if !websocket.Client.Active() {
		go websocket.Client.Start(config.Values.Port)
	}

	Initialize()
	go StartScanner()

	<-*websocket.Client.OnReadySub()

	ctnfile.WriteFile(delPath, []string{"bout to be gone"})
	ctnfile.WriteFile(newPath, []string{"brand new"})
	ctnfile.WriteFile(newExPath, []string{"new but excluded"})
	ctnfile.WriteFile(modPath, []string{"has been modified"})

	time.Sleep(time.Second * 1)

	err = os.Remove(delPath)
	if err != nil {
		clog.Error(err)
	}

	time.Sleep(time.Second * 4)
}
