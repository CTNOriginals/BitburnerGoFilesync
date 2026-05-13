package commands

import (
	"fmt"
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

		fmt.Printf("%s\n", args[0].Def.StringRecurse(filter...))
	},
}

func GetFilePathAutoComplete(root *string) func(string) []string {
	return func(line string) []string {
		var parts = strings.Split(line, " ")
		var partial = parts[len(parts)-1]

		var searchDir string
		var dirPrefix string

		if strings.HasPrefix(partial, "/") {
			var idx = strings.LastIndex(partial, "/")
			if idx >= 0 {
				searchDir = partial[:idx]
				dirPrefix = partial[:idx+1]
			}
			if searchDir == "" {
				searchDir = "/"
			}
		} else {
			var idx = strings.LastIndex(partial, "/")
			if idx >= 0 {
				searchDir = filepath.Join(*root, partial[:idx])
				dirPrefix = partial[:idx+1]
			} else {
				searchDir = *root
				dirPrefix = ""
			}
		}

		var entries, err = os.ReadDir(searchDir)
		if err != nil {
			return nil
		}

		var names = make([]string, 0, len(entries))
		for _, entry := range entries {
			var name = entry.Name()
			if entry.IsDir() {
				names = append(names, dirPrefix+name+"/")
			} else {
				names = append(names, dirPrefix+name)
			}
		}

		return names
	}
}
