package cli

import (
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
	"github.com/chzyer/readline"
)

type Option struct {
	Name        string
	Description []string

	Callback readline.DynamicCompleteFunc

	Children OptionList
}

func (this Option) IsDynamic() bool {
	return this.Callback != nil
}

func (this Option) Build() *readline.PrefixCompleter {
	var opt = readline.PrefixCompleter{
		Name:     []rune(this.Name),
		Dynamic:  this.IsDynamic(),
		Callback: this.Callback,
		Children: make([]readline.PrefixCompleterInterface, len(this.Children)),
	}

	for i, child := range this.Children {
		opt.Children[i] = child.Build()
	}

	return &opt
}

func (this Option) String() string {
	var str strings.Builder

	str.WriteString("-")

	if this.IsDynamic() {
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

	if len(this.Children) > 0 {
		str.WriteString("\n")
		str.WriteString(ctnstring.Indent(this.Children.String(), 1, "  "))
	}

	return str.String()
}
