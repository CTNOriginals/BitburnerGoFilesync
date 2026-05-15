package cmdconfig

import (
	"fmt"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

var def = commands.Definition{
	Name: "config",
	Description: []string{
		"List all config fields and their values.",
	},

	Options: commands.TList{
		// set_def, // TODO:
	},
	Execution: func(args commands.TInputList, self int) {
		fmt.Printf("%s: %v\n", "Port", config.Values.Port)
		fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
		fmt.Printf("%s: %v\n", "ConfigFilePath", constants.ConfigFilePath)
		fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
		fmt.Printf("%s: %v\n", "FileScanInterval", config.Values.FileScanInterval)
		fmt.Printf("%s: [%v]\n", "IncludeFileExt", strings.Join(config.Values.FilePatterns.Include, ", "))
		fmt.Printf("%s: [%v]\n", "ExcludeFileExt", strings.Join(config.Values.FilePatterns.Exclude, ", "))
	},
}

func init() {
	commands.List.Push(&def)
}
