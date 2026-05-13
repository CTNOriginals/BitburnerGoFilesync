package commands

import (
	"fmt"
	"slices"
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

	// Prevent this Definition from being included in any string
	Hidden bool

	AutoComplete readline.DynamicCompleteFunc
	Validator    FnValidator
	Execution    FnExecution
}

func (this *Definition) ValidateSelf() {
	for _, opt := range this.Options {
		opt.ValidateSelf()
	}

	if this.Name == HelperSubCommand.Name {
		return
	}

	if this.Options.GetDefinitionByName(HelperSubCommand.Name) == nil {
		this.Options.Push(HelperSubCommand)
	}
}

func (this Definition) HasAutoComplete() bool {
	return this.AutoComplete != nil
}

func (this Definition) Build() readline.PrefixCompleterInterface {
	var build = &Custom_PrefixCompleter{
		PrefixCompleter: readline.PrefixCompleter{
			Name:     []rune(this.Name),
			Dynamic:  this.HasAutoComplete(),
			Callback: this.AutoComplete,
			Children: make([]readline.PrefixCompleterInterface, len(this.Options)),
		},
	}

	if !this.HasAutoComplete() && this.ExpectValue {
		build.Dynamic = true
		build.Callback = autocompleteVoid
	}

	for i, opt := range this.Options {
		build.Children[i] = opt.Build()
	}

	return build
}

func (this Definition) stringDescription() string {
	var str strings.Builder
	var width = len(this.Name)

	if len(this.Description) > 0 {
		str.WriteString(": ")
		str.WriteString(this.Description[0])
	}

	if len(this.Description) > 1 {
		str.WriteString("\n")
		var mid = strings.Join(this.Description[1:len(this.Description)-1], "\n")
		var bot = this.Description[len(this.Description)-1]

		mid = ctnstring.Indent(mid, 1, "├ ")
		bot = ctnstring.Indent(bot, 1, "└ ")

		var desc = ""
		if len(this.Description) > 2 {
			desc = fmt.Sprintf("%s\n%s", mid, bot)
		} else {
			desc = fmt.Sprintf("%s", bot)
		}

		desc = ctnstring.Indent(desc, width, " ")

		str.WriteString(desc)
	}

	return str.String()
}

func (this Definition) stringHead() string {
	var str strings.Builder

	str.WriteString(this.Name)

	var desc = this.stringDescription()
	str.WriteString(desc)

	return str.String()
}

func (this Definition) String() string {
	return this.stringHead()
}

func (this Definition) StringRecurse(filter ...string) string {
	var str strings.Builder

	if len(filter) > 0 && !slices.Contains(filter, this.Name) {
		return str.String()
	}

	var lines = strings.Split(this.stringHead(), "\n")
	str.WriteString(lines[0])

	var optionString = this.Options.StringRecurse(filter...)
	if len(optionString) > 0 {
		optionString = ctnstring.Indent(optionString, 2, " ")
		lines = append(lines, strings.Split(optionString, "\n")...)
	}

	var chars = []rune("┃┣┗")
	var optionCount = len(this.Options) - 1

	for i := 1; i < len(lines); i++ {
		str.WriteString("\n")

		if optionCount <= 0 {
			str.WriteString(lines[i])
			continue
		}

		var line = []rune(lines[i])
		var char = chars[0]

		if !strings.ContainsRune(string(chars)+" ", line[2]) {
			if optionCount > 1 {
				char = chars[1]
			} else {
				char = chars[2]
			}

			line[1] = '╸'
			optionCount -= 1
		}

		line[0] = char
		str.WriteString(string(line))
	}

	return str.String()
}
