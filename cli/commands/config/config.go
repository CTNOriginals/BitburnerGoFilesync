package cmdconfig

import (
	"log"
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
		log.Printf("%s: %v\n", "Port", config.Values.Port)
		log.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
		log.Printf("%s: %v\n", "ConfigFilePath", constants.ConfigFilePath)
		log.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
		log.Printf("%s: %v\n", "FileScanInterval", config.Values.FileScanInterval)
		log.Printf("%s: [%v]\n", "IncludeFileExt", strings.Join(config.Values.FilePatterns.Include, ", "))
		log.Printf("%s: [%v]\n", "ExcludeFileExt", strings.Join(config.Values.FilePatterns.Exclude, ", "))
	},
}

func init() {
	commands.List.Push(&def)
}
