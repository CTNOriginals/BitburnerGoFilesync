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
	var options = make([]readline.PrefixCompleterInterface, len(this))

	for i, def := range this {
		var opts = def.BuildOptions()
		opts.Name = []rune(def.Name)

		options[i] = opts
	}

	return readline.NewPrefixCompleter(options...)
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
