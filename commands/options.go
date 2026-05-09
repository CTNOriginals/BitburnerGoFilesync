package commands

import "github.com/chzyer/readline"

type Options struct {
	Triggers    []string
	Description []string

	Callback readline.DynamicCompleteFunc
}

func (this Options) IsDynamic() bool {
	return this.Callback == nil
}
