package cmdconfig

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

var def = commands.Definition{
	Name: "config",
	Description: []string{
		"Config interface, with multiple functions:",
		"1. view the config",
		"2. edit config values",
	},

	Options: commands.TList{
		{Name: "list",
			Description: []string{"Lists all config fields and values."},
		},
		{Name: "set",
			Description: []string{
				"Set a config fields value.",
				"If the new value should be remembered, pass --save.",
				"some more args here...",
			},

			Options: commands.TList{
				{Name: "field",
					Description:  []string{"The config field to set."},
					AutoComplete: getField,
					Options: commands.TList{
						{Name: "value",
							Description:  []string{"The value to set."},
							AutoComplete: getField,
						},
					},
				},
			},
		},
	},
	Execution: execute,
}

func execute(args commands.TInputList) {
	fmt.Printf("%s: %v\n", "Port", config.Values.Port)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	fmt.Printf("%s: %v\n", "ConfigDirectory", constants.ConfigFilePath)
	fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
	fmt.Printf("%s: %v\n", "IncludeFileExt", config.Values.FilePatterns.Include)
	fmt.Printf("%s: %v\n", "FileScanDelay", config.Values.FilePatterns.Exclude)
	fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
}

func getField(line string) []string {
	fmt.Printf("\nconfig set getField: %s\n", line)

	return []string{}
}

func init() {
	commands.List = append(commands.List, &def)
}
