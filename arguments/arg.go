package arguments

import (
	"fmt"
	"slices"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

type argDef struct {
	Alias       []string
	Description []string

	Params argParameters
	Action func(params []string)
}

func (this argDef) String() string {
	var alias = strings.Join(this.Alias, ", ")
	var desc = strings.Join(this.Description, "\n"+ctnstring.Repeat(" ", 4))
	var params = ctnstring.Indent(this.Params.String(), 2, " ")

	var str = alias

	if len(this.Description) > 0 {
		str += fmt.Sprintf(":\n    %s", desc)
	}

	if len(this.Params) > 0 {
		str += fmt.Sprintf("\n%s", params)
	}

	return str + "\n"
}

type argList []*argDef

func (this argList) GetDefByAlias(alias string) (def *argDef, exists bool) {
	for _, arg := range this {
		if slices.Contains(arg.Alias, alias) {
			return arg, true
		}
	}

	return
}
