package arguments

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/test"

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
			"If an argument doesnt have its required parameters, it will say so in the console,",
			"this argument will not execute anything after that and will be ignored.",
			"If you pass in more parameters than an argument needs, it simply ignores the overflow.",
		},
	},

	{Alias: []string{"--help", "--wtf"},
		Description: []string{
			"Prints a list of arguments and their descriptions.",
			"Follow it up with another argument (without the -- before it)",
			"to get a more detailed explanation about that argument.",
		},
		Params: argParameters{
			{
				Name: "command",
				Description: []string{
					"The name (without the -- before it) of a command.",
					"Print a detailed explanation about a specific command.",
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

			os.Exit(0)
		},
	},
	{Alias: []string{"--full-help", "--fhelp"},
		Description: []string{
			"The same as --help, but it also includes all of the extra information",
			"as if you entered --help <command> for each argument.",
		},
		Params: argParameters{},
		Action: func(params []string) {
			for _, arg := range onInitList {
				println(arg.String())
			}

			os.Exit(0)
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
				os.Exit(1)
			}

			if !ctnfile.PathExists(params[0]) {
				fmt.Printf("'--dir' directory does not exist: %s\n", params[0])
				os.Exit(1)
			}

			constants.SetBitburnerDir(params[0])
		},
	},
	{Alias: []string{"--include-ext", "--ext"},
		Description: []string{
			"Specify which file extensions the file watcher should include.",
			"Default: js ts txt",
		},
		Params: argParameters{
			{Name: "extensions", Description: []string{
				"Any number of file extensions separated with spaces.",
				"Example: js ts json",
			}},
		},
		Action: func(params []string) {
			constants.IncludeFileExt = params
		},
	},
	{Alias: []string{"--port"},
		Description: []string{
			"Set the port for the server to connect to.",
			"By default, the server will try to connect to 'localhost:8080'.",
		},
		Params: argParameters{
			{Name: "port", Description: []string{
				"The port number.",
				"Default: 8080",
			}},
		},
		Action: func(params []string) {
			if len(params) == 0 {
				fmt.Print("'--port' requires at least 1 parameter.\n")
				os.Exit(1)
			}

			constants.Port = params[0]
		},
	},
	{Alias: []string{"--scan-interval", "--interval"},
		Description: []string{
			"The amount of miliseconds the file scanner waits each loop.",
			"By default 100, if <= 0 it will skip the sleep function entirely.",
		},
		Params: argParameters{
			{Name: "interval", Description: []string{
				"The interval in miliseconds",
				"Default: 100",
			}},
		},
		Action: func(params []string) {
			if len(params) == 0 {
				fmt.Println("'--scan-interval' requires at least 1 parameter.")
				os.Exit(1)
			}

			if !ctnstring.Validate(params[0], "1234567890") {
				fmt.Println("'--scan-interval' only accepts number characters")
				os.Exit(1)
			}

			var num, err = strconv.ParseInt(params[0], 0, 10)

			if err != nil {
				fmt.Println("'--scan-interval'", err)
				os.Exit(1)
			}

			constants.FileScanDelay = int(num)
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
			os.Exit(0)
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
				"Usually when a server is ran, the program wont exit as it keeps evaluating it,",
				"if this parameter is set to true, the program will still be prevented from exiting.",
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
