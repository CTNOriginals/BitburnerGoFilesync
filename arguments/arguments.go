package arguments

import (
	"runtime"
	"strconv"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/test"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

// NOTE: This list can be used inside the actions of arguments
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
			if len(params) == 0 {
				printHelp()
			} else {
				printHelpSelect(params...)
			}

			runtime.Goexit()
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
				clog.Message(arg.String())
			}

			runtime.Goexit()
		},
	},
	{Alias: []string{"--config"},
		Description: []string{
			"Define the config.toml file path.",
			"By default, the config file is located in the same directory as the binary.",
			"If no config file exists at the specified location, one will be created.",
		},
		Params: argParameters{
			{Name: "filepath", Description: []string{
				"The path to the config file",
				"Default: " + constants.ConfigFilePath,
			}},
		},
		Action: func(params []string) {
			if len(params) == 0 {
				clog.Error("'--config' requires at least 1 parameter.\n")
				runtime.Goexit()
			}

			constants.ConfigFilePath = params[0]
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
				clog.Error("'--dir' requires at least 1 parameter.\n")
				runtime.Goexit()
			}

			if !ctnfile.PathExists(params[0]) {
				clog.Errorf("'--dir' directory does not exist: %s\n", params[0])
				runtime.Goexit()
			}

			config.ValidateBitburnerDirectory(params[0])
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
			config.Values.FilePatterns.Include = params
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
				clog.Error("'--port' requires at least 1 parameter.\n")
				runtime.Goexit()
			}

			config.Values.Port = params[0]
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
				clog.Error("'--scan-interval' requires at least 1 parameter.")
				runtime.Goexit()
			}

			if !ctnstring.Validate(params[0], "1234567890") {
				clog.Error("'--scan-interval' only accepts number characters")
				runtime.Goexit()
			}

			var num, err = strconv.ParseInt(params[0], 0, 64)

			if err != nil {
				clog.Error("'--scan-interval'", err)
				runtime.Goexit()
			}

			config.Values.FileScanInterval = int(num)
		},
	},
	{Alias: []string{"--get-definitions"},
		Description: []string{
			"Currently not functional.",
			"Requests the NetscriptDefinitions.d.ts file when a connection is established.",
			"The definitions file will be created in bitburners root directory.",
		},
		Params: argParameters{},
		Action: func(params []string) {
			clog.Error("TODO: Handle --get-definitions")
		},
	},

	{Alias: []string{"DEBUG ARGUMENTS"}},

	{Alias: []string{"--test", "--debug"},
		Description: []string{
			"Runs the test function if it exists",
		},
		Params: argParameters{},
		Action: func(params []string) {
			constants.Debug = true
			test.DoTest(params...)
			runtime.Goexit()
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
	{Alias: []string{"--no-server", "--no-client", "--no-websocket"},
		Description: []string{
			"Prevents the program from creating a server and connecting to bitburner.",
		},
		Params: argParameters{},
		Action: func(params []string) {
			constants.NoServer = true
		},
	},
	{Alias: []string{"--no-cli"},
		Description: []string{
			"Prevents the program running the cli.",
		},
		Params: argParameters{},
		Action: func(params []string) {
			constants.NoServer = true
		},
	},
}
