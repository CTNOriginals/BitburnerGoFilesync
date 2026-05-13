package commands

import (
	"fmt"
	"os"
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

func FileComplete(filePath *string) func(string) []string {
	return func(line string) []string {
		// var parts = strings.Split(line, " ")
		var names = make([]string, 0)
		var files, _ = os.ReadDir(*filePath)

		// fmt.Printf("\n%s/%s\n", *filePath, parts[len(parts)-1])

		for _, file := range files {
			names = append(names, file.Name())
		}

		return names
	}
}
