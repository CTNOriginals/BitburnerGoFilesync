package entries

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/commands"
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
}
