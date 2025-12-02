package main

import (
	"filesync/communication"
	"filesync/constants"
	"filesync/debug"
	"filesync/test"
	"filesync/watcher"
	"fmt"
	"os"
	"time"
)

var (
	debugMode = false
	testMode  = false
	noServer  = false
	noWatcher = false
	keepAlive = false
)

func main() {
	parseArgs(os.Args)
	// fmt.Printf("debugMode: %t\ntestMode: %t\nnoServer: %t\nnoWatcher: %t\nkeepAlive: %t", debugMode, testMode, noServer, noWatcher, keepAlive)

	startTime := time.Now()
	fmt.Printf("\n\n---- FileSync START %s ----\n", startTime.Format(time.TimeOnly))
	defer fmt.Printf("---- FileSync END %s ----\n", startTime.Format(time.TimeOnly))

	if testMode {
		test.DoTest()
	}

	if debugMode {
		go debug.DebugCommandListener()
	}

	if !noWatcher {
		watcher.Initialize()
		go watcher.FileScanner()
	}

	if !noServer {
		communication.StartServer(constants.Port)
	} else if keepAlive {
		for {
			time.Sleep(time.Millisecond)
		}
	}
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch arg {
		case "--test":
			testMode = true
		case "--debug":
			debugMode = true
		case "--no-watcher":
			noWatcher = true
		case "--no-server":
			noServer = true
		case "--keep-alive":
			keepAlive = true
		}
	}
}
