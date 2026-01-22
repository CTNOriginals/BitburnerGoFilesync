package test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
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

	filePatternMatching()
}

func filePatternMatching() {
	// watcher.Initialize()

	config.Values.FilePatterns.Include = []string{
		// "[",
		"*.ts*",
		"*.js*",
	}
	config.Values.FilePatterns.Exclude = []string{
		"*.d.ts",
		"test.ts",
	}
	var testPaths = []string{
		"test.ts",
		"path/index.ts",
		"main.tsx",
		// "some.thingts",
		"some.things.jswtf",
		// ".git/objects",
		"path/to/file.js",
		"path/file.d.ts",
		"path/file.d.ts.ts",
	}

	var match = func(pattern string, path string, strict bool) bool {
		var variants = []string{
			"**/" + pattern,
			"*/**/" + pattern,
		}

		var valid bool = utils.Expect(filepath.Match(pattern, path))
		// fmt.Printf("  %s: %t\n", pattern, valid)

		for _, variant := range variants {
			var state = utils.Expect(filepath.Match(variant, path))
			// fmt.Printf("  %s: %t\n", variant, state)

			if strict {
				valid = valid && state
			} else {
				valid = valid || state
			}
		}

		return valid
	}

	var accepted = []string{}

	for _, path := range testPaths {
		fmt.Printf("%s\n", path)
		var valid bool = true

		println(" Exclude:")

		for _, pattern := range config.Values.FilePatterns.Exclude {
			fmt.Printf("  %s: %t\n", pattern, match(pattern, path, false))

			if !match(pattern, path, false) {
				continue
			}

			valid = false
			break
		}

		if !valid {
			println()
			continue
		}

		println(" Include:")
		valid = false
		for _, pattern := range config.Values.FilePatterns.Include {
			fmt.Printf("  %s: %t\n", pattern, match(pattern, path, false))

			if !match(pattern, path, false) {
				continue
			}

			valid = true
			break
		}

		println()
		if !valid {
			continue
		}

		accepted = append(accepted, path)
	}

	fmt.Printf("Accepted:\n  %s\n\n", strings.Join(accepted, ",\n  "))

	// fmt.Printf("%v\n", strings.Join(includes, ",\n"))
	// println("\n")
	// fmt.Printf("%v\n", strings.Join(excludes, ",\n"))

	// for _, path := range testPaths {
	// 	var state = watcher.ShouldIncludeFile(path)
	// 	// fmt.Printf("%t: %s\n", state, path)
	// 	_ = state
	// }
}
