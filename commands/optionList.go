package commands

import "strings"

type OptionList []Option

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
