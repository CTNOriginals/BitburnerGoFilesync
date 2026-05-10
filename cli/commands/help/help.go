package cmdhelp

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
)

var def = commands.Definition{
	Triggers: []string{"help", "h", "?"},
	Description: []string{
		"Print a list of all possible commands",
		"along with their descriptions.",
		"Follow it with the name of any specific command",
		"for more info on just that command.",
	},

	Options: commands.OptionList{
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
			Callback: dynamic_getCommands,
		},
	},

	Execution: execute,
}

func dynamic_getCommands(line string) []string {
	// TODO:
	return []string{}
}

func execute(args ...string) {
	fmt.Printf("%s\n", commands.List)
}

func init() {
	commands.List = append(commands.List, &def)
}
