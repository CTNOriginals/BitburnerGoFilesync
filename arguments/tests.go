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

	definitions()
	findDefinitions()
	parseSafeArgs()
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
		"--no-server", "true",
		"--include-ext", "js", "ts",
	})

	clog.Infof("  NoWatcher: %t\n", constants.NoWatcher)
	clog.Infof("  NoServer:  %t\n", constants.NoServer)
	clog.Infof("  KeepAlive: %t\n", constants.KeepAlive)
	clog.Infof("  Include:   %v\n", config.Values.FilePatterns.Include)
}
