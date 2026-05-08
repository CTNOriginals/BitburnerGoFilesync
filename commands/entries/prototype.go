package entries

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
	"github.com/chzyer/readline"
)

var proto = commands.Definition{
	Triggers: []string{"prototype", "proto"},
	Description: []string{
		"A prototype command used for testing",
	},

	Options: readline.NewPrefixCompleter(
		readline.PcItem("wah"),
		commands.ReadLine_FileItem,
	),

	Execution: proto_execute,
}

func proto_execute(args ...string) {
	fmt.Printf("proto waah!\n%v\n", args)
}

func init() {
	commands.CommandList = append(commands.CommandList, &proto)
}
