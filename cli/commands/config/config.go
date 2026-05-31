package cmdconfig

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

var clog = clogger.Default

var def = commands.Definition{
	Name: "config",
	Description: []string{
		"List all config fields and their values.",
	},

	Options: commands.TList{
		// set_def, // TODO:
	},
	Execution: func(args commands.TInputList, self int) {
		clog.Messagef("%s: %v\n", "Port", config.Values.Port)
		clog.Messagef("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
		clog.Messagef("%s: %v\n", "ConfigFilePath", constants.ConfigFilePath)
		clog.Messagef("%s: %v\n", "BitburnerRoot", config.Values.Directory)
		clog.Messagef("%s: %v\n", "FileScanInterval", config.Values.FileScanInterval)
		clog.Messagef("%s: [%v]\n", "IncludeFileExt", strings.Join(config.Values.FilePatterns.Include, ", "))
		clog.Messagef("%s: [%v]\n", "ExcludeFileExt", strings.Join(config.Values.FilePatterns.Exclude, ", "))
	},
}

func init() {
	commands.List.Push(&def)
}
