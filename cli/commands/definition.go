package commands

import (
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
	"github.com/chzyer/readline"
)

type Definition struct {
	// The strings that will activate this command
	Name        string
	Description []string

	Options OptionList

	Execution func(args ...string)
}

func (this Definition) BuildOptions() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(this.Options.Build()...)
}

func (this Definition) String() string {
	var str = strings.Builder{}

	str.WriteString(this.Name)

	str.WriteString(":\n  ")
	str.WriteString(strings.Join(this.Description, "\n  "))

	str.WriteString("\n")
	str.WriteString(ctnstring.Indent(this.Options.String(), 1, "  "))

	return str.String()
}
