package main

import (
	"filesync/communication"
	"filesync/constants"
	"filesync/debug"
	"filesync/test"
	"fmt"
	"os"
	"time"
)

var (
	debugMode = false
	testMode  = false
)

func main() {
	parseArgs(os.Args)

	startTime := time.Now()
	fmt.Printf("\n\n---- FileSync START %s ----\n", startTime.Format(time.TimeOnly))
	defer fmt.Printf("---- FileSync END %s ----\n", startTime.Format(time.TimeOnly))

	if testMode {
		test.DoTest()
		return
	}

	if debugMode {
		go debug.DebugCommandListener()
	}

	communication.StartServer(constants.Port)
}

func parseArgs(args []string) {
	for _, arg := range args {
		switch arg {
		case "--test":
			testMode = true
		case "--debug":
			debugMode = true
		}
	}
}
