package commands

import (
	"strings"

	"github.com/chzyer/readline"
)

type OptionList []Option

func (this OptionList) Build() []readline.PrefixCompleterInterface {
	var options = make([]readline.PrefixCompleterInterface, len(this))

	for i, opt := range this {
		options[i] = opt.Build()
	}

	return options
}

func (this OptionList) String() string {
	var str strings.Builder

	for _, opt := range this {
		if str.Len() > 0 {
			str.WriteString("\n")
		}

		str.WriteString(opt.String())
	}

	return str.String()
}
