package commands

import "strings"

type TList []*Definition

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
