package arguments

import (
	"slices"
)

type argStream struct {
	Alias  string
	Params []string
}

// Parse specific args removes any args and their params
// from args if they are not included in list
// or if include is false, the args in list will be remove from args instead.
func ParseSpecificArgs(args []string, include bool, list ...string) {
	if len(list) == 0 {
		if !include {
			ParseArgs(args)
		}

		return
	}

	var validArgs = []string{}
	var validArg = false

	for _, part := range args {
		if len(part) >= 2 && part[0:2] == "--" {
			validArg = include == slices.Contains(list, part)
		}

		if validArg {
			validArgs = append(validArgs, part)
		}
	}

	ParseArgs(validArgs)
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
		stream = append(stream, argStream{Alias: currentKey, Params: currentStream})
	}

	for i, part := range args {
		if len(part) >= 2 && part[0:2] == "--" {
			appendStream()
			currentKey = part
			currentParams = currentParams[:0]
		} else {
			currentParams = append(currentParams, part)
		}

		if i == len(args)-1 {
			appendStream()
		}
	}

	for _, stream := range stream {
		var def, exists = argumentList.GetDefByAlias(stream.Alias)

		if !exists {
			clog.Errorf("Unknown argument flag: %s\n", stream.Alias)
			continue
		}

		def.Action(stream.Params)
	}
}
