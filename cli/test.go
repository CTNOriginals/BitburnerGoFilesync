package cli

import (
	"fmt"
	"strings"
)

func TestCli() {
	// fmt.Printf("%s\n", commands.List.String())

	// cliSegments()
	// cliCommands()
	// cliHelp()

	CommandWatcher()
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

		var segments, err = GetInputSegments(line)

		if err == nil {
			fmt.Printf("%s\n", strings.Join(segments, "\n"))
		} else {
			fmt.Printf("%v\n", err)
		}
	}
}

func cliTester(commands []string) {
	for _, line := range commands {
		fmt.Printf(">> %s\n", line)

		var err = ParseInput(line)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		println("")
	}
}

func cliCommands() {
	fmt.Printf("\n-- Command Tests --\n")

	cliTester([]string{
		"help",
		"help config",
		"help full",
		"prototype \"foo bar\"",
		"prototype 1242e+8",
		"prototype logging/index.ts",
		"config set Port 1234",
		"config list",
	})
}

func cliHelp() {
	fmt.Printf("\n-- Help Tests --\n")

	cliTester([]string{
		"prototype help",
		"prototype num help",
		"prototype 'foo bar' help",
		"prototype logging/index.ts help",
		"config set help",
		"config set Port 1234 help",
	})
}
