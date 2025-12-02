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
	var desc = strings.Join(this.Description, "\n"+ctnstring.Repeat(" ", len(alias)+2))

	var params = ctnstring.Indent(this.Params.String(), len(alias), " ")
	// params = strings.Join(strings.Split(params, "\n"), "\n  ")

	return fmt.Sprintf("%s: %s\n%s\n", alias, desc, params)
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
