package cmdhelp

import (
	"fmt"
	"slices"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

var def = commands.Definition{
	Name: "help",
	Description: []string{
		"Print a list of all possible commands",
		"along with their descriptions.",
		"Follow it with the name of any specific command",
		"for more info on just that command.",
	},

	Options: commands.TList{
		{Name: "full",
			Description: []string{
				"Print all of the enformation about each command",
				"instead of only the name and description.",
			},
		},
		{Name: "command",
			Description: []string{
				"Print one commands info exclusively",
				"and potentially with more info.",
			},
			AutoComplete: dynamic_getCommands,
			Validator: func(input string) bool {
				return slices.Contains(commands.List.GetNamesRecursive(), input)
			},
		},
	},

	Execution: execute,
}

func dynamic_getCommands(line string) []string {
	var triggers = make([]string, len(commands.List))

	for i, def := range commands.List {
		triggers[i] = def.Name
	}

	return triggers
}

func execute(args commands.TInputList) {
	fmt.Printf("%s\n", commands.List)
}

func init() {
	commands.List = append(commands.List, &def)
}
