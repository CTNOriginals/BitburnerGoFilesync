package cmdhelp

import (
	"fmt"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

var def = commands.Definition{
	Name: "help",
	Description: []string{
		"Print a list of all possible commands",
		"along with their descriptions and sub commands.",
	},

	Execution: func(args commands.TInputList, self int) {
		var instructions = []string{
			"",
			"-- COMMAND INSTRUCTIONS --",
			"Help & Info:",
			" To get more information about any (sub)command,",
			" simply type out the commands you want more info about",
			" and append the word 'help' at the end of the input line.",
			" This will print the entire tree of those commands and their options.",
			" Example:",
			"  config set help - prints config>set>all subcommands of 'set'",
			"  config set Port 1234 help - prints config>set>field>value",
			"",
			"Tab Completion:",
			" Tab completion works differently depending on your current input.",
			" Most of the time, a Tab press just lists the currently available commands.",
			" If your current input is a command that expects a value instead of a sub command name,",
			" the Completion may list some possible values, but sometimes this is not possible.",
		}
		fmt.Printf("%s\n\n%s\n", strings.Join(instructions, "\n"), commands.List.StringRecurse())
	},
}

func init() {
	commands.List.Push(&def)
}
