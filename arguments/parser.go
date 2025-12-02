package arguments

import "fmt"

type argStream struct {
	Alias  string
	Params []string
}

func ParseArgs(args []string) {
	onInitList = argumentList

	var stream []argStream

	var currentKey = ""
	var currentParams []string

	var appendStream = func() {
		if currentKey == "" {
			return
		}

		var currentStream = make([]string, len(currentParams))
		copy(currentStream, currentParams)
		stream = append(stream, argStream{Alias: currentKey, Params: currentParams})
	}

	for i, part := range args {
		if part[0:2] == "--" {
			appendStream()
			currentKey = part
			currentParams = currentParams[:0]
			continue
		}

		currentParams = append(currentParams, part)

		if i == len(args)-1 {
			appendStream()
		}
	}

	for _, stream := range stream {
		var def, exists = argumentList.GetDefByAlias(stream.Alias)

		if !exists {
			fmt.Printf("Unknown argument flag: %s\n", stream.Alias)
			continue
		}

		def.Action(stream.Params)
	}
}
