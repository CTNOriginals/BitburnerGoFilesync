package cli

import (
	"fmt"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

func TestCli() {
	clog.Infof("Commands:\n%s\n", commands.List.StringRecurse())

	testSegments()
	// cliCommands()
	// cliHelp()

	CommandWatcher()
}

func cliTester(commands []string) {
	for _, line := range commands {
		clog.Infof(">> %s\n", line)

		var err = ParseInput(line)
		if err != nil {
			clog.Errorf("%v\n", err)
		}
	}
}

func cliCommands() {
	clog.Message("-- Command Tests --")

	cliTester([]string{
		// "config",
		"prototype 123",
		"prototype \"42\"",
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

func wrapSegments(segments []string) string {
	var builder strings.Builder

	for _, seg := range segments {
		if builder.Len() > 0 {
			builder.WriteRune(' ')
		}

		fmt.Fprintf(&builder, "(%s)", seg)
	}

	return builder.String()
}

func testSegments() {
	clog.Message("\n-- Segment Tests --")

	var tests = []string{
		"prototype \"42\"",
		"prototype 'io x' x",
		"config set \"foo bar\" baz",
		"\"foo 'baz' bar\" goo 'drap bah'",
		"\"foo 'baz",
		"something 'like a dog' or cat'",
		"some weird='edge ca'se` that i hope` ne've`r happens",
		"goo foo='bar ins\\'t baz' but is",
	}

	for _, line := range tests {
		var segments = GetInputSegments(line)
		clog.Infof("\n%s\n%s", line, wrapSegments(segments))
	}
}
