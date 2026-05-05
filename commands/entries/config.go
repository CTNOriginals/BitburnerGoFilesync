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

	Execution: config_execute,
}

func init() {
	commands.CommandList = append(commands.CommandList, &conf)
}

func config_execute() {
	fmt.Printf("%s: %v\n", "Port", config.Values.Port)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.ConfigFilePath)
	fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
	fmt.Printf("%s: %v\n", "IncludeFileExt", config.Values.FilePatterns.Include)
	fmt.Printf("%s: %v\n", "FileScanDelay", config.Values.FilePatterns.Exclude)
	fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
}
