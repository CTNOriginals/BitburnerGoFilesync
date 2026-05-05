package entries

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
)

var def = commands.Definition{
	Triggers: []string{"prototype", "proto"},
	Description: []string{
		"A prototype command used for testing",
	},

	Execution: execute,
}

func init() {
	commands.CommandList = append(commands.CommandList, &def)
}

func execute() {
	fmt.Printf("proto waah!\n")
}
