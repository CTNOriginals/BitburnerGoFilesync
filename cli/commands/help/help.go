package cmdhelp

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

var def = commands.Definition{
	Name: "help",
	Description: []string{
		"Print a list of all possible commands",
		"along with their descriptions.",
		"For more detailed info on any command",
		"type out what ever command and any options",
		"and append 'help' at the end of it to print the",
		"info specific to those commands and options.",
		// TODO: add more clear instructions on how commands are used
		// make those instructions seperate from this help command.
	},

	Options: commands.TList{},

	Execution: func(args commands.TInputList, self int) {
		fmt.Printf("%s\n", commands.List.StringRecurse())
	},
}

func init() {
	commands.List.Push(&def)
}
