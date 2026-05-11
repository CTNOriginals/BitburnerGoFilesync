package commands

import (
	"strings"

	"github.com/chzyer/readline"
)

type TList []*Definition

func (this TList) GetCommandByTrigger(trigger string) *Definition {
	for _, def := range this {
		if def.Name == trigger {
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
