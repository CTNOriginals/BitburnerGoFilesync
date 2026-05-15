package commands

import "fmt"

type Input struct {
	Def   *Definition
	Value string
}

func (this Input) String() string {
	return fmt.Sprintf("%s: %s", this.Def.Name, this.Value)
}
