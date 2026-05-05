package commands

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var CommandList []*Definition = make([]*Definition, 0)
var promtSymbol = "\033[34m\033[1m»\033[0m "

func CommandWatcher() {
	// fmt.Printf("commands: %v\n", CommandList)
	fmt.Print(promtSymbol)

	var scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var cmd = strings.TrimSpace(scanner.Text())

		var def = GetCommandByTrigger(cmd)

		if def == nil {
			fmt.Printf("Unknown command: %s\n", cmd)
			goto endscan
		}

		def.Execution()

	endscan:
		fmt.Print(promtSymbol)
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
