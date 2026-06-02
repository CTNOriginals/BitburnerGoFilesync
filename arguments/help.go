package arguments

import (
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

func printHelp() {
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

		clog.Messagef("%s: %s\n\n", alias, desc)
	}
}

func printHelpSelect(aliases ...string) {
	for _, alias := range aliases {
		var def, exists = onInitList.GetDefByAlias("--" + alias)
		if !exists {
			clog.Errorf("Unknown argument flag: %s\n", alias)
			continue
		}

		clog.Message(def.String())
	}

}
