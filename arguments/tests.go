package arguments

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/test"
)

func init() {
	test.Register("arguments", TestArgs)
}

func TestArgs() {
	clog.Info("\n-- Arguments Tests --\n")

	// definitions()
	// findDefinitions()
	// parseSafeArgs()
	// parseStream()

	printHelp()
}

func definitions() {
	clog.Message("\n-- Definitions --\n")

	for _, def := range argumentList {
		if len(def.Alias) == 0 || !strings.HasPrefix(def.Alias[0], "--") {
			continue
		}

		clog.Messagef("%s\n", def.String())
	}
}

func parseStream() {
	var tests = []string{
		"--unkown foo bar",
		// "--help",
		// "--port",
		// "--scan-interval",
		// "--wtf wtf", // exits after
	}

	for _, stream := range tests {
		clog.Infof("%s\n", stream)
		ParseArgs(strings.Split(stream, " "))
	}
}

func findDefinitions() {
	clog.Message("\n-- Find Definitions --\n")

	ParseArgs([]string{})

	var tests = []string{
		"--help",
		"--config",
		"--dir",
		"--unknown-flag",
	}

	for _, alias := range tests {
		var def, exists = onInitList.GetDefByAlias(alias)
		if exists {
			clog.Infof("  %s -> %s\n", alias, strings.Join(def.Alias, ", "))
		} else {
			clog.Infof("  %s -> not found\n", alias)
		}
	}
}

func parseSafeArgs() {
	clog.Message("\n-- Parse Safe Args --\n")

	ParseArgs([]string{
		"--no-watcher",
		"--no-server",
		"--no-cli",
		"--include-ext", "js", "ts", "test",
	})

	clog.Infof("  NoWatcher: %t\n", constants.NoWatcher)
	clog.Infof("  NoServer:  %t\n", constants.NoServer)
	clog.Infof("  NoCli:     %t\n", constants.NoCli)
	clog.Infof("  Include:   %v\n", config.Values.FilePatterns.Include)
}
