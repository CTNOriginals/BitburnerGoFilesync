package test

import (
	"fmt"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	// "github.com/CTNOriginals/BitburnerGoFilesync/config"
	// "github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func DoTest() {
	// config.Initialize()
	// println("\n")
	// fmt.Printf("%s: %v\n", "Port", config.Values.Port)
	// fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	// fmt.Printf("%s: %v\n", "WorkindDirectory", constants.ConfigFilePath)
	// fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
	// fmt.Printf("%s: %v\n", "IncludeFileExt", config.Values.FilePatterns.Include)
	// fmt.Printf("%s: %v\n", "FileScanDelay", config.Values.FilePatterns.Exclude)
	// fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	// fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	// fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
	// println("")
	// watcher.FileScanner()
	TestCli()
}

func TestCli() {
	fmt.Printf("%s\n", commands.List)

	// cliSegments()
	// cliCommands()

	cli.CommandWatcher()
	// readlineDemo()
}

func cliSegments() {
	fmt.Printf("\n-- Segment Tests --\n")

	var tests = []string{
		"config set \"foo bar\" baz",
		"\"foo 'baz' bar\" goo 'drap bah'",
		"\"foo 'baz",
		"something 'like a dog' or cat'",
		"some weird='edge ca'se` that i hope` ne've`r happens",
		"goo foo='bar ins\\'t baz' but is",
	}

	for _, line := range tests {
		fmt.Printf(">> %s\n", line)

		var segments, err = cli.GetInputSegments(line)

		if err == nil {
			fmt.Printf("%s\n", strings.Join(segments, "\n"))
		} else {
			fmt.Printf("%v\n", err)
		}
	}
}

func cliCommands() {
	fmt.Printf("\n-- Command Tests --\n")

	var commandTests = []string{
		"help",
		"help config",
		"help full",
		"prototype \"foo bar\"",
		"prototype 1242e+8",
		"prototype logging/index.ts",
		"config set Port 1234",
		"config list",
	}

	for _, line := range commandTests {
		fmt.Printf(">> %s\n", line)

		var err = cli.ParseInput(line)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		println("")
	}
}
