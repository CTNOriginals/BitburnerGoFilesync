package test

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/watcher"
)

func DoTest() {
	config.Initialize()
	println("\n")
	fmt.Printf("%s: %v\n", "Port", config.Values.Port)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
	fmt.Printf("%s: %v\n", "IncludeFileExt", config.Values.FilePatterns.Include)
	fmt.Printf("%s: %v\n", "FileScanDelay", config.Values.FilePatterns.Exclude)
	fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
	println("")

}

func filePatternMatching() {
	// watcher.Initialize()

	config.Values.FilePatterns.Include = []string{
		"*.ts",
		"*.js*",
	}
	config.Values.FilePatterns.Exclude = []string{
		"*.d.ts",
		"to/*",
		"test.ts",
	}
	var testPaths = []string{
		"test.ts",
		"some.thingts",
		"some.things.jswtf",
		".git/objects",
		"path/to/file.js",
		"path/file.d.ts",
		"path/file.d.ts.ts",
	}

	for _, path := range testPaths {
		var state = watcher.ShouldIncludeFile(path)
		// fmt.Printf("%t: %s\n", state, path)
		_ = state
	}
}
