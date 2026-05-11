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

	Options TList

	AutoComplete readline.DynamicCompleteFunc
	Validator    func(input string) bool
	Execution    func(args ...string)
}

func (this Definition) HasAutoComplete() bool {
	return this.AutoComplete != nil
}

func (this Definition) Build() *readline.PrefixCompleter {
	var build = readline.PrefixCompleter{
		Name:     []rune(this.Name),
		Dynamic:  this.HasAutoComplete(),
		Callback: this.AutoComplete,
		Children: make([]readline.PrefixCompleterInterface, len(this.Options)),
	}

	for i, opt := range this.Options {
		build.Children[i] = opt.Build()
	}

	return &build
}

func (this Definition) String() string {
	var str strings.Builder

	str.WriteString("-")

	if this.HasAutoComplete() {
		str.WriteRune('@')
	} else {
		str.WriteRune(' ')
	}

	str.WriteString(this.Name)

	var width = str.Len()

	if len(this.Description) > 0 {
		str.WriteString(": ")
		str.WriteString(this.Description[0])
	}

	if len(this.Description) > 1 {
		str.WriteString("\n")
		var desc = strings.Join(this.Description[1:], "\n")
		desc = ctnstring.Indent(desc, 1, "| ")
		str.WriteString(ctnstring.Indent(desc, width, " "))
	}

	if len(this.Options) > 0 {
		str.WriteString("\n")
		str.WriteString(ctnstring.Indent(this.Options.String(), 1, "  "))
	}

	return str.String()
}
