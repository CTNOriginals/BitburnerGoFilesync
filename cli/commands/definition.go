package commands

import (
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
	"github.com/chzyer/readline"
)

type FnValidator func(input string) bool
type FnExecution func(args TInputList, self int)

var autocompleteVoid readline.DynamicCompleteFunc = func(s string) []string { return nil }

// Definition contains all of the info that can build and execute a command
type Definition struct {
	Name        string
	Description []string

	Options TList

	// Wether to expect a value instead of the Name of the command
	// same as adding a validator that always returns true.
	// This will also prevent the Name from showing up as an autocompleted suggestion.
	// This field does not need to be set if the Validator is defined.
	ExpectValue bool

	AutoComplete readline.DynamicCompleteFunc
	Validator    FnValidator
	Execution    FnExecution
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

	if !this.HasAutoComplete() && this.ExpectValue {
		build.Dynamic = true
		build.Callback = autocompleteVoid
	}

	for i, opt := range this.Options {
		build.Children[i] = opt.Build()
	}

	return &build
}

func (this Definition) stringHead() string {
	var str strings.Builder

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

	return str.String()
}

func (this Definition) String() string {
	var str strings.Builder

	str.WriteString(this.stringHead())

	// for _, opt := range this.Options {
	// 	str.WriteString("\n")
	// 	str.WriteString(ctnstring.Indent(opt.stringHead(), 2, " "))
	// }

	return str.String()
}

func (this Definition) StringRecurse() string {
	var str strings.Builder

	str.WriteString("-")

	if this.HasAutoComplete() {
		str.WriteRune('@')
	} else {
		str.WriteRune(' ')
	}

	str.WriteString(this.stringHead())

	if len(this.Options) > 0 {
		str.WriteString("\n")
		str.WriteString(ctnstring.Indent(this.Options.String(), 1, "  "))
	}

	return str.String()
}
