package commands

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var CommandList []*Definition = make([]*Definition, 0)

func CommandWatcher() {
	fmt.Printf("commands: %v\n", CommandList)
	fmt.Printf("User Input: ")

	var scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var cmd = strings.TrimSpace(scanner.Text())

		fmt.Printf("Received: %s\n", cmd)

		var def = GetCommandByTrigger(cmd)

		if def == nil {
			fmt.Printf("Unknown command: %s\n", cmd)
			goto endscan
		}

		def.Execution()

	endscan:
		fmt.Printf("User Input: ")
	}
}

func GetCommandByTrigger(trigger string) *Definition {
	for _, def := range CommandList {
		if slices.Contains(def.Triggers, trigger) {
			return def
		}
	}

	return nil
}
