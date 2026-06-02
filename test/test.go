package test

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/watcher"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var testFunctions = map[string]func(){
	"logger":  clogger.TestLogger,
	"client":  websocket.TestClient,
	"cli":     cli.TestCli,
	"watcher": watcher.TestWatcher,
}

func Register(name string, fn func()) {
	testFunctions[name] = fn
}

func DoTest(args ...string) {
	log.Printf("Running debug with args: %v\n", args)

	for _, arg := range args {
		var fn, exists = testFunctions[arg]

		if !exists {
			log.Printf("Unknown test function name: %s\n", arg)
			continue
		}

		go fn()
	}
}
