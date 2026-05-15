package commands

import (
	"strings"

	"github.com/chzyer/readline"
)

type Custom_PrefixCompleter struct {
	readline.PrefixCompleter
}

func (c *Custom_PrefixCompleter) GetDynamicNames(line []rune) [][]rune {
	var names = [][]rune{}
	for _, name := range c.Callback(string(line)) {
		if strings.HasSuffix(name, "/") {
			names = append(names, []rune(name))
		} else {
			names = append(names, []rune(name+" "))
		}
	}
	return names
}
