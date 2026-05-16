package cmdproto

import (
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

var def = commands.Definition{
	Name: "prototype",
	Description: []string{
		"A prototype command used for testing",
	},

	Options: commands.TList{
		{Name: "nodesc",
			Description: []string{},
		},
		{Name: "file",
			Description:  []string{"A file path for something."},
			AutoComplete: commands.GetFilePathAutoComplete(&config.Values.Directory),
			Validator: func(input string) bool {
				return (path.IsAbs(input) && ctnfile.FileExists(input)) ||
					ctnfile.FileExists(utils.GetAbsolutePath(input))
			},
		},
		{Name: "str",
			Description: []string{"any string"},
			Validator: func(input string) bool {
				var quotes = "\"'`"
				var l = input[0]
				var r = input[len(input)-1]

				return strings.ContainsRune(quotes, rune(l)) &&
					strings.ContainsRune(quotes, rune(r)) &&
					l == r
			},
		},
		{Name: "num",
			Description: []string{"any number"},
			Validator: func(input string) bool {
				var _, floaterr = strconv.ParseFloat(input, 64)
				return floaterr == nil
			},
		},
	},

	Execution: func(args commands.TInputList, self int) {
		log.Printf("%v\n", args)
	},
}

func init() {
	commands.List.Push(&def)
}
