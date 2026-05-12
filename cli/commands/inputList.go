package commands

import (
	"fmt"
	"strings"
)

type TInputList []*Input

func (this TInputList) Execute() {
	fmt.Printf("%s\n", this)
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
