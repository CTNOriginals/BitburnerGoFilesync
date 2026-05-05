package commands

import (
	"fmt"
	"strings"
)

type Definition struct {
	// The strings that will activate this command
	Triggers    []string
	Description []string

	Execution func()
}

func (this Definition) String() string {
	var str = strings.Builder{}

	str.WriteString(this.Triggers[0])

	if len(this.Triggers) > 1 {
		str.WriteString(fmt.Sprintf(" (%s)", strings.Join(this.Triggers[1:], ", ")))
	}

	str.WriteString(":\n  ")
	str.WriteString(strings.Join(this.Description, "\n  "))

	return str.String()
}
