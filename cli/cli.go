package cli

import (
	"fmt"
	"io"
	"log"
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
		var inputs, inputErr = commands.List.ParseInput(parts, nil)

		if inputErr != nil {
			fmt.Printf("Input error: %s\n", inputErr)
			continue
		}

		inputs.Execute()
	}
}
