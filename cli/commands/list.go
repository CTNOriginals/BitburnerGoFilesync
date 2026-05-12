package commands

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

type TList []*Definition

func (this TList) GetDefinitionByName(name string) *Definition {
	for _, def := range this {
		if def.Name == name {
			return def
		}
	}

	return nil
}

func (this TList) GetNames() []string {
	var names = make([]string, len(this))

	for i, def := range this {
		names[i] = def.Name
	}

	return names
}

func (this TList) GetNamesRecursive() []string {
	var names = this.GetNames()

	for _, def := range this {
		names = append(names, def.Options.GetNamesRecursive()...)
	}

	return names
}

func (this TList) Build() *readline.PrefixCompleter {
	var build = make([]readline.PrefixCompleterInterface, len(this))

	for i, def := range this {
		build[i] = def.Build()
	}

	return readline.NewPrefixCompleter(build...)
}

// Runs the input through each definitions validator until one returns true
func (this TList) TryGetValidatedDefinition(input string) *Definition {
	for _, def := range this {
		if def.Validator == nil || !def.Validator(input) {
			continue
		}

		return def
	}

	// if no defined validator succeded,
	// check if any autocomplete suggested this input
	for _, def := range this {
		if def.ExpectValue ||
			// BUG: this may produce inconsistency if the autocomplete
			// function uses any of the preceding line of inputs
			(def.HasAutoComplete() && slices.Contains(def.AutoComplete(""), input)) {
			return def
		}
	}

	return nil
}

func (this TList) ParseInput(input string, args ...string) (TInputList, error) {
	var def = this.GetDefinitionByName(input)

	if def == nil {
		// fmt.Printf("%v > %s (%v): def nil\n", this.GetNames(), input, args)
		def = this.TryGetValidatedDefinition(input)

		if def == nil {
			// fmt.Printf("%s (%v): def nil\n", input, args)
			return nil, fmt.Errorf("Unknown input: %s\n", input)
		}
	}

	var inputList = TInputList{&Input{
		Def:   def,
		Value: input,
	}}

	if len(args) > 0 {
		var argList, err = def.Options.ParseInput(args[0], args[1:]...)

		if err != nil {
			return nil, err
		}

		inputList = append(inputList, argList...)
	}

	return inputList, nil
}

func (this TList) string(recurse bool) string {
	var str strings.Builder

	for _, def := range this {
		if str.Len() > 0 {
			str.WriteRune('\n')
		}

		if recurse {
			str.WriteString(def.StringRecurse())
		} else {
			str.WriteString(def.String())
		}
	}

	return str.String()
}

func (this TList) String() string {
	return this.string(false)
}
func (this TList) StringRecurse() string {
	return this.string(true)
}
