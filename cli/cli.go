package cli

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/chzyer/readline"
)

var promtSymbol = "\033[34m\033[1m»\033[0m "

func CommandWatcher() {
	var options = commands.List.Build()

	var cli, err = readline.NewEx(&readline.Config{
		Prompt:       promtSymbol,
		AutoComplete: options,
		// HistoryFile:       "/tmp/readline.tmp", // TODO: support other os's
		HistorySearchFold:   true,
		ForceUseInteractive: true,
	})

	if err != nil {
		panic(fmt.Sprintf("Commands readline error:\n%v", err))
	}

	defer cli.Close()

	log.SetOutput(cli.Stderr())

	for {
		var line, err = cli.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		var parts = strings.Split(line, " ")
		var prefix = parts[0]
		var args = parts[1:]
		var found = false

		if slices.Contains([]string{"help", "?"}, prefix) {
			fmt.Printf("%s\n", commands.List.String())
			continue
		}

		for _, def := range commands.List {
			if !def.IsTrigger(prefix) {
				continue
			}

			found = true
			def.Execution(args...)
			break
		}

		if !found {
			fmt.Printf("Unknown command: %s\n", prefix)
		}
	}
}
