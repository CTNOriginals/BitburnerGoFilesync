package main

import (
	"filesync/constants"
	"filesync/debug"
	"filesync/rpc"
	"filesync/test"
	"os"
)

var (
	debugMode = false
	testMode  = false
)

func main() {
	parseArgs(os.Args)

	print("\n\n---- FileSync START ----\n")
	defer print("---- Program END ----")

	if testMode {
		test.DoTest()
		return
	}

	if debugMode {
		go debug.DebugCommandListener()
	}

	rpc.StartServer(constants.Port)
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
