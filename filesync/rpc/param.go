package rpc

import "fmt"

func (this Param) String() string {
	str := ""

	for i, param := range this {
		str += fmt.Sprintf("\n\t\"%s\": \"%s\"", param[0], param[1])

		if i != len(this)-1 {
			str += ","
		}
	}

	if str != "" {
		str += "\n"
	}

	return str
}

func (this *Param) Add(key string, val string) {
	*this = append(*this, [2]string{key, val})
}
