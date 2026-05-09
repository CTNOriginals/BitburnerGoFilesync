package commands

import "github.com/chzyer/readline"

type Option struct {
	Name        string
	Description []string

	Callback readline.DynamicCompleteFunc

	Children []Option
}

func (this Option) IsDynamic() bool {
	return this.Callback == nil
}

func (this Option) Build() *readline.PrefixCompleter {
	var opt = readline.PrefixCompleter{
		Name:     []rune(this.Name),
		Dynamic:  this.IsDynamic(),
		Callback: this.Callback,
		Children: make([]readline.PrefixCompleterInterface, len(this.Children)),
	}

	for i, child := range this.Children {
		opt.Children[i] = child.Build()
	}

	return &opt
}
