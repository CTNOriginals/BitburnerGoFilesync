package arguments

import (
	"filesync/constants"
	"filesync/test"
	"fmt"
	"strings"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

// This list can be used inside the actions of arguments
// to workaround the initialization cycle error.
// It will be assigned once ParseArgs is called.
var onInitList argList = nil

var argumentList = argList{
	{Alias: []string{"Formatting Rules"},
		Description: []string{
			"Each new argument always has to start with a double dash '--'.",
			"If the argument does not start with '--' it is considered a parameter",
			"for the most recent argument that started with '--'.",
			"",
			"Each argument may have any number of parameters,",
			"to check what an argument may accept or require,",
			"you can do --help followed by the name of the argument without the '--'.",
			"",
			"Some arguments may accept a specific amount of parameters where others accept a range.",
			"If an argument doesnt have its required parameters, it will say say so in the console,",
			"this argument will not execute anything after that and will be ignored.",
			"If you pass in more parameters then an argument needs, it simply ignores the overflow.",
		},
	},

	{Alias: []string{"--help", "--wtf"},
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
			if len(params) == 0 {
				fmt.Print("'--dir' requires at least 1 parameter.\n")
				return
			}

			if !ctnfile.PathExists(params[0]) {
				fmt.Printf("Directory does not exist: %s\n", params[0])
				return
			}

			constants.SetBitburnerDir(params[0])
		},
	},
	{Alias: []string{"--include-ext", "--ext"},
		Description: []string{
			"Specify which file extensions the file watcher should include.",
			"If ommited, all file types will be included inside the root --dir",
		},
		Params: argParameters{
			{Name: "extensions", Description: []string{
				"Any number of file extensions seperated with spaces.",
				"Example: js ts json",
			}},
		},
		Action: func(params []string) {
			constants.IncludeFileExt = params
		},
	},

	{Alias: []string{"DEBUG ARGUMENTS"}},

	{Alias: []string{"--test"},
		Description: []string{
			"Runs the test function if it exists",
		},
		Params: argParameters{},
		Action: func(params []string) {
			test.DoTest()
		},
	},
	{Alias: []string{"--no-watcher"},
		Description: []string{
			"Prevents the program from watching file events.",
		},
		Params: argParameters{},
		Action: func(params []string) {
			constants.NoWatcher = true
		},
	},
	{Alias: []string{"--no-server"},
		Description: []string{
			"Prevents the program from creating a server and connecting to bitburner.",
		},
		Params: argParameters{
			{Name: "keep-alive", Description: []string{
				"Accepts: true, false",
				"Usually when a server is ran, the program wont exit because as it keeps evaluating it,",
				"if this param is set to true, the program will still be prevented from exiting without a server active.",
			}},
		},
		Action: func(params []string) {
			constants.NoServer = true

			if len(params) > 0 && params[0] == "true" {
				constants.KeepAlive = true
			}
		},
	},
	// {Alias: []string{"--config"},
	// 	Description: []string{
	// 		"COMING SOON: Specify the location of the config file.",
	// 	},
	// 	Params: argParameters{
	// 		{Name: "file", Description: []string{
	// 			"The file path to the config",
	// 		}},
	// 	},
	// 	Action: func(params []string) {
	// 		println("TODO: Add a global config")
	// 	},
	// },
}
