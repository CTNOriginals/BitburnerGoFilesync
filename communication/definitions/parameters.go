package definitions

import (
	"fmt"
)

// The parameters of a request are always of type string.
// This type will only hold the fields of any params that need to be passed
type ParameterFields []string

func (this ParameterFields) Generate(values []string) string {
	if len(this) != len(values) {
		fmt.Printf("ParameterFields  (%v) do not match passed in values (%v)", this, values)
		return "{}"
	}

	str := ""

	for i, field := range this {
		str += fmt.Sprintf("  \"%s\": \"%s\"", field, values[i])

		if i < len(this)-1 {
			str += ",\n"
		}
	}

	return fmt.Sprintf("{\n%s\n }", str)
}
