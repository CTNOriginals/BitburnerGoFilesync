package cmdproto

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var def = commands.Definition{
	Name: "prototype",
	Description: []string{
		"A prototype command used for testing",
	},

	Options: commands.OptionList{
		{Name: "wah",
			Description: []string{},
		},
		{Name: "file",
			Description:  []string{"A file path for something."},
			AutoComplete: cli.Readline_FileComplete(&config.Values.Directory),
		},
	},

	Execution: execute,
}

func execute(args ...string) {
	fmt.Printf("proto waah!\n%v\n", args)
}

func init() {
	commands.List = append(commands.List, &def)
}
