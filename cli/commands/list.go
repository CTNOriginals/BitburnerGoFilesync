package commands

import (
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

type TList []*Definition

func (this TList) GetCommandByTrigger(trigger string) *Definition {
	for _, def := range this {
		if slices.Contains(def.Triggers, trigger) {
			return def
		}
	}

	return nil
}

func (this TList) Build() *readline.PrefixCompleter {
	var options = make([]readline.PrefixCompleterInterface, len(this))

	for i, def := range this {
		var opts = def.BuildOptions()
		opts.Name = []rune(def.Triggers[0])

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
