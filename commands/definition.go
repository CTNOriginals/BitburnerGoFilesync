package commands

import (
	"fmt"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
	"github.com/chzyer/readline"
)

type Definition struct {
	// The strings that will activate this command
	Triggers    []string
	Description []string

	Options OptionList

	Execution func(args ...string)
}

func (this Definition) IsTrigger(compare string) bool {
	compare = strings.TrimSpace(compare)
	compare = strings.ToLower(compare)

	for _, trigger := range this.Triggers {
		trigger = strings.ToLower(trigger)

		if compare == trigger {
			return true
		}
	}

	return false
}

func (this Definition) BuildOptions() *readline.PrefixCompleter {
	var options = make([]readline.PrefixCompleterInterface, len(this.Options))

	for i, opt := range this.Options {
		options[i] = opt.Build()
	}

	return readline.NewPrefixCompleter(options...)
}

func (this Definition) String() string {
	var str = strings.Builder{}

	str.WriteString(this.Triggers[0])

	if len(this.Triggers) > 1 {
		str.WriteString(fmt.Sprintf(" (%s)", strings.Join(this.Triggers[1:], ", ")))
	}

	str.WriteString(":\n  ")
	str.WriteString(strings.Join(this.Description, "\n  "))

	str.WriteString("\n")
	str.WriteString(ctnstring.Indent(this.Options.String(), 1, "  "))

	return str.String()
}
