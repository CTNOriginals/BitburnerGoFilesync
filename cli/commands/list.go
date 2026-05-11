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

func (this TList) ParseInput(line []string, parent *Input) (TInputList, error) {
	var input = line[0]
	var def = this.GetDefinitionByName(input)

	if def == nil {
		return nil, fmt.Errorf("Unknown input: %s\n", input)
	}

	var parsed = TInputList{&Input{
		Def:   def,
		Value: input,
	}}

	if len(line) > 1 {
	}

	return parsed, nil
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
