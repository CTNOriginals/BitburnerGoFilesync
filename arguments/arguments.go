package arguments

import (
	"fmt"
	"strings"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

// This list is only used by the help argument action
// to workaround the initialization cycle error
var onInitList argList = nil

var argumentList = argList{
	{Alias: []string{"--help", "--wtf", "?"},
		Description: []string{
			"Prints a list of arguments and their descriptions.",
			"Follow it up with another argument (without the -- before it)",
			"to get a more detailed explination about that argument.",
		},
		Params: argParameters{
			{
				Name: "command",
				Description: []string{
					"The name (without the -- before it) of a command.",
					"Print a detailed explination about a specific command.",
				},
			},
		},
		Action: func(params []string) {
			if len(params) > 0 {
				for _, param := range params {
					var def, exists = onInitList.GetDefByAlias("--" + params[0])
					if !exists {
						fmt.Printf("Unknown argument flag: %s\n", param)
						continue
					}

					println(def.String())
				}

				return
			}

			var maxAliasSpace = 0

			// Precalculate the max amount of spaces any alias will ever take in
			// to then be able to apply that space before the descrition of each argument
			for _, def := range onInitList {
				var length = len(strings.Join(def.Alias, ", "))
				if length > maxAliasSpace {
					maxAliasSpace = length
				}
			}

			for _, def := range onInitList {
				var alias = strings.Join(def.Alias, ", ")
				var desc = ctnstring.Repeat(" ", maxAliasSpace-len(alias))
				desc += strings.Join(def.Description, "\n"+ctnstring.Repeat(" ", maxAliasSpace+2))

				fmt.Printf("%s: %s\n\n", alias, desc)
			}

			println("")
		},
	},
	{Alias: []string{"--config"},
		Description: []string{
			"Specify the location of the config file.",
		},
		Params: argParameters{
			{Name: "file", Description: []string{
				"The file path to the config",
			}},
		},
		Action: func(params []string) {
			println("TODO: Link the config file")
		},
	},
	{Alias: []string{"--dir"},
		Description: []string{
			"Specify the directory where this tool should watch",
			"for file changes to sync up with bitburner",
		},
		Params: argParameters{
			{Name: "dir", Description: []string{
				"The path to the directory where you keep your bitburner scripts.",
				"Make sure to surround this parameter with double quotes (\").",
			}},
		},
		Action: func(params []string) {
			println("TODO: Set the directory path")
			if !ctnfile.PathExists(params[0]) {
				fmt.Printf("Directory does not exist: %s\n", params[0])
				return
			}
		},
	},
}
