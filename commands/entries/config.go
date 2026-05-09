package entries

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

var conf = commands.Definition{
	Triggers: []string{"config", "conf", "info", "settings", "options"},
	Description: []string{
		"List all config fields along with their current values.",
	},

	Options: []commands.Option{
		{Name: "list",
			Description: []string{"Lists all config fields and values."},
		},
		{Name: "set",
			Description: []string{"Lists all config fields and values."},

			Children: []commands.Option{
				{Name: "field",
					Description: []string{"The config field to set."},
					Callback:    config_getField,
				},
				{Name: "value",
					Description: []string{"The value to set."},
					Callback:    config_getField,
				},
			},
		},
	},
	Execution: config_execute,
}

func config_execute(args ...string) {
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

func config_getField(line string) []string {
	fmt.Printf("\nconfig set getField: %s\n", line)

	return []string{}
}

func init() {
	commands.CommandList = append(commands.CommandList, &conf)
}
