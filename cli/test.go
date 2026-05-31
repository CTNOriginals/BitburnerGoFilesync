package cli

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

func TestCli() {
	clog.Infof("Commands:\n%s\n", commands.List.StringRecurse())

	cliSegments()
	// cliCommands()
	// cliHelp()

	CommandWatcher()
}

func cliSegments() {
	clog.Info("\n-- Segment Tests --\n")

	var tests = []string{
		"config set \"foo bar\" baz",
		"\"foo 'baz' bar\" goo 'drap bah'",
		// "\"foo 'baz",
		"something 'like a dog' or cat'",
		// "some weird='edge ca'se` that i hope` ne've`r happens",
		"goo foo='bar ins\\'t baz' but is",
	}

	for _, line := range tests {
		clog.Infof(">> %s\n", line)

		var segments, err = GetInputSegments(line)

		if err == nil {
			clog.Messagef("%s\n", strings.Join(segments, "\n"))
		} else {
			clog.Errorf("%v\n", err)
		}
	}
}

func cliTester(commands []string) {
	for _, line := range commands {
		clog.Infof(">> %s\n", line)

		var err = ParseInput(line)
		if err != nil {
			clog.Infof("%v\n", err)
		}
	}
}

func cliCommands() {
	clog.Message("\n-- Command Tests --\n")

	cliTester([]string{
		"config",
		"prototype 123",
	})
}

func cliHelp() {
	clog.Message("\n-- Help Tests --\n")

	cliTester([]string{
		"prototype help",
		"prototype num help",
		"prototype 'foo bar' help",
		"prototype foo.js help",
		// "config set help",
		// "config set Port 1234 help",
	})
}
