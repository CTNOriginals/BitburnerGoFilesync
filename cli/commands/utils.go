package commands

import (
	"os"
	"path/filepath"
	"strings"
)

var HelperSubCommand = &Definition{
	Name:        "help",
	Description: []string{"Print the info about all commands and options infront of this."},
	Hidden:      true,

	AutoComplete: func(s string) []string {
		return nil
	},
	Validator: func(input string) bool {
		input = strings.TrimSpace(input)
		return input == "help"
	},
	Execution: func(args TInputList, self int) {
		var filter = make([]string, len(args)-1)

		for i, arg := range args {
			if i == len(args)-1 {
				break
			}

			filter[i] = arg.Def.Name
		}

		// list all posibilities downwards
		filter = append(filter, args[len(args)-2].Def.Options.GetNamesRecursive()...)

		clog.Messagef("%s\n", args[0].Def.StringRecurse(filter...))
	},
}

func GetFilePathAutoComplete(root *string) func(string) []string {
	return func(line string) []string {
		var parts = strings.Split(line, " ")
		var partial = parts[len(parts)-1]

		var dir, prefix string
		if idx := strings.LastIndex(partial, "/"); idx >= 0 {
			dir = partial[:idx]
			prefix = partial[:idx+1]
		}

		var searchDir = *root
		switch {
		case strings.HasPrefix(partial, "/"):
			searchDir = dir
			if searchDir == "" {
				searchDir = "/"
			}
		case dir != "":
			searchDir = filepath.Join(*root, dir)
		}

		var entries, err = os.ReadDir(searchDir)
		if err != nil {
			return nil
		}

		var names = make([]string, 0, len(entries))
		for _, e := range entries {
			var name = e.Name()
			if e.IsDir() {
				name += "/"
			}
			names = append(names, prefix+name)
		}

		return names
	}
}
