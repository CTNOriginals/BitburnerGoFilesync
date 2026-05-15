package commands

import (
	"slices"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
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

	if !this.HasAutoComplete() {
		if this.ExpectValue {
			// NOTE: Prevent the Name from showing up as an autocomplete option
			build.Dynamic = true
			build.Callback = autocompleteVoid
		} else {
			// NOTE: This causes the name to be autocompleted with a space at the end.
			// It serves as a QoL feature so that the user can keep typing right after completion.
			build.Name = append(build.Name, ' ')
		}
	}

	for i, opt := range this.Options {
		build.Children[i] = opt.Build()
	}

	return build
}

func (this Definition) StringDescription() string {
	var str strings.Builder
	var width = len(this.Name)

	if len(this.Description) > 0 {
		str.WriteString(": ")
		str.WriteString(this.Description[0])
	}

	if len(this.Description) > 1 {
		var desc = strings.Join(this.Description, "\n")

		desc = ctnstring.Indent(desc, 2, " ")
		desc = utils.FormatStringAsTree(&strings.Builder{}, utils.TreeFormatOptions{
			IndentCount: 2,
			Whitespace:  []rune(" "),
			MaxDepth:    len(this.Description) - 1,
			Thin:        true,
		}, strings.Split(desc, "\n")...)

		desc = ctnstring.Indent(desc, width, " ")

		str.WriteString(desc)
	}

	return str.String()
}

func (this Definition) String() string {
	var str strings.Builder

	str.WriteString(this.Name)

	var desc = this.StringDescription()
	str.WriteString(desc)

	return str.String()
}

func (this Definition) StringRecurse(filter ...string) string {
	var str strings.Builder

	if len(filter) > 0 && !slices.Contains(filter, this.Name) {
		return str.String()
	}

	var lines = strings.Split(this.String(), "\n")
	str.WriteString(lines[0])

	var optionString = this.Options.StringRecurse(filter...)
	if len(optionString) > 0 {
		optionString = ctnstring.Indent(optionString, 2, " ")
		lines = append(lines, strings.Split(optionString, "\n")...)
	}

	utils.FormatStringAsTree(&str, utils.TreeFormatOptions{
		IndentCount: 2,
		Whitespace:  []rune(" "),
		MaxDepth:    len(this.Options) - 1,
	}, lines...)

	return str.String()
}
