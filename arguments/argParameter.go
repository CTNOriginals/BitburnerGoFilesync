package arguments

import (
	"fmt"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

type argParamDef struct {
	Name        string
	Description []string
}

type argParameters []argParamDef

func (this argParameters) String() string {
	var str = ""

	if len(this) > 0 {
		str = "Parameters:\n"
	}

	for i, param := range this {
		var desc = strings.Join(param.Description, "\n"+ctnstring.Repeat(" ", 4))

		str += fmt.Sprintf("  %s:\n    %s", param.Name, desc)

		if i < len(this)-1 {
			str += "\n"
		}
	}

	return str
}
