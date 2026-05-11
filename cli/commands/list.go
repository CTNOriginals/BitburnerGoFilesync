package commands

import (
	"fmt"
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

	return nil
}

func (this TList) ParseInput(input string, args ...string) (TInputList, error) {
	var def = this.GetDefinitionByName(input)

	if def == nil {
		def = this.TryGetValidatedDefinition(input)

		if def == nil {
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

func (this TList) String() string {
	var str strings.Builder

	for _, def := range this {
		if str.Len() > 0 {
			str.WriteRune('\n')
		}

		str.WriteString(def.String())
	}

	return str.String()
}
