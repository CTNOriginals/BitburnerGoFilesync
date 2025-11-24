package definitions

import "fmt"

// The parameters of a request are always of type string.
// This type will only hold the fields of any params that need to be passed
type ParameterFields []string

func (this ParameterFields) Generate(values []string) string {
	str := ""

	for i, field := range this {
		if i > len(values)-1 {
			fmt.Printf("There are not enough values to populate the remaining parameter fields: %v", this[i:])
			break
		}

		str += fmt.Sprintf("\"%s\": \"%s\"", field, values[i])

		if i < len(this)-1 {
			str += ",\n"
		}
	}

	return fmt.Sprintf("{\n  %s\n }", str)
}
