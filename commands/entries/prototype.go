package entries

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
)

var proto = commands.Definition{
	Triggers: []string{"prototype", "proto"},
	Description: []string{
		"A prototype command used for testing",
	},

	Execution: proto_execute,
}

func init() {
	commands.CommandList = append(commands.CommandList, &proto)
}

func proto_execute() {
	fmt.Printf("proto waah!\n")
}
