package arguments

import "strings"

// The max amount of space any single alias will ever take.
var maxAliasWidth = -1

func getMaxAliasWidth() int {
	var size = 0

	for _, def := range onInitList {
		for _, alias := range def.Alias {
			if len(alias) > size {
				size = len(alias)
			}
		}
	}

	return size
}

func getSingleHelp(def argDef) string {
	var builder strings.Builder

	// the current alias width
	var currentWidth = 0
	var line = 0

	var newLine = func() {
		if line <= len(def.Description)-1 {
			builder.WriteString(strings.Repeat(" ", (maxAliasWidth-currentWidth)+2))
			builder.WriteString(def.Description[line])
		}

		currentWidth = 0
		line += 1
		builder.WriteRune('\n')
	}

	for _, alias := range def.Alias {
		if currentWidth+len(alias) > maxAliasWidth {
			newLine()
		}

		builder.WriteString(alias)
		builder.WriteRune(' ')
		currentWidth += len(alias) + 1
	}

	for line <= len(def.Description)-1 {
		newLine()
	}

	var block = builder.String()

	// trim any leading newline char
	// the return should only be the content
	// not ending in a newline to prevent confusion
	block = strings.TrimSuffix(block, "\n")

	return block
}

func printHelp() {
	if maxAliasWidth == -1 {
		maxAliasWidth = getMaxAliasWidth()
	}

	var builder strings.Builder

	for _, def := range onInitList {
		if builder.Len() > 0 {
			builder.WriteString("\n\n")
		}

		builder.WriteString(getSingleHelp(*def))
	}

	clog.Message(builder.String())
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
