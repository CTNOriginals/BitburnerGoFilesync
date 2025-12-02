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

	for _, param := range this {
		if str != "" {
			str += "\n"
		}

		var desc = strings.Join(param.Description, "\n"+ctnstring.Repeat(" ", len(param.Name)+4))

		str += fmt.Sprintf("- %s: %s", param.Name, desc)
	}

	return str
}
