package main

import (
	"log"
	"os"
	"slices"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/arguments"
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/watcher"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "main",
})

func main() {
	log.SetFlags(0)

	var startTime = time.Now()
	clog.Messagef("\n\n---- FileSync START %s ----\n", startTime.Format(time.TimeOnly))
	defer func() {
		var now = time.Now()
		clog.Messagef(
			"---- FileSync END %s (%s) ----\n",
			now.Format(time.TimeOnly),
			time.Since(startTime).String(),
		)

		// Make sure that all gorotines also terminate when runtime.Goexit is called
		os.Exit(0)
	}()

	var args = os.Args

	constants.Debug = slices.Contains(args, "--test") || slices.Contains(args, "--debug")

	// Make sure that the config file path is the correct one before initializing the config
	arguments.ParseSpecificArgs(args, true, "--config")
	config.Initialize()

	arguments.ParseSpecificArgs(args, false, "--config")

	if !constants.NoCli {
		go cli.CommandWatcher()
	}
	if !constants.NoWatcher {
		watcher.Initialize()
		go watcher.FileScanner()
	}

	if !constants.NoServer {
		go websocket.Client.Start(config.Values.Port)
		defer websocket.Client.Close()
	}

	for {
		time.Sleep(time.Second)
	}
}
