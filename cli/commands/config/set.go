package cmdconfig

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

var set_def = &commands.Definition{
	Name: "set",
	Description: []string{
		"Set a config fields value.",
		"If the new value should be remembered, pass --save.",
	},

	Options: commands.TList{
		{Name: "field",
			Description:  []string{"The config field to set."},
			AutoComplete: getField,
			Options: commands.TList{
				{Name: "value",
					Description: []string{"The value to set."},
					ExpectValue: true,
				},
			},
		},
	},
}

func getField(line string) []string {
	return ctnstruct.Keys(config.Values)
}
