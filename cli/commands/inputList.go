package commands

import (
	"fmt"
	"strings"
)

type TInputList []*Input

func (this TInputList) Execute() {
	for i := len(this) - 1; i >= 0; i-- {
		var input = this[i]

		if input.Def.Execution == nil {
			continue
		}

		input.Def.Execution(this, i)
		break
	}
}

func (this TInputList) String() string {
	var str strings.Builder

	for i, input := range this {
		if str.Len() > 0 {
			str.WriteRune('\n')
		}

		fmt.Fprintf(&str, "%d: %s", i, input)
	}

	return str.String()
}
