package cmdproto

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var proto = commands.Definition{
	Triggers: []string{"prototype", "proto"},
	Description: []string{
		"A prototype command used for testing",
	},

	Options: commands.OptionList{
		{Name: "wah",
			Description: []string{},
		},
		{Name: "file",
			Description: []string{"A file path for something."},
			Callback:    commands.Readline_FileComplete(&config.Values.Directory),
		},
	},

	Execution: proto_execute,
}

func proto_execute(args ...string) {
	fmt.Printf("proto waah!\n%v\n", args)
}

func init() {
	commands.CommandList = append(commands.CommandList, &proto)
}
