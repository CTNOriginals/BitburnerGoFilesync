package cmdproto

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var proto = cli.Definition{
	Triggers: []string{"prototype", "proto"},
	Description: []string{
		"A prototype command used for testing",
	},

	Options: cli.OptionList{
		{Name: "wah",
			Description: []string{},
		},
		{Name: "file",
			Description: []string{"A file path for something."},
			Callback:    cli.Readline_FileComplete(&config.Values.Directory),
		},
	},

	Execution: proto_execute,
}

func proto_execute(args ...string) {
	fmt.Printf("proto waah!\n%v\n", args)
}

func init() {
	cli.CommandList = append(cli.CommandList, &proto)
}
