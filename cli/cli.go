package cli

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/chzyer/readline"
)

func CommandWatcher() {
	var cli, err = readline.NewEx(&readline.Config{
		Prompt:       "\033[34m\033[1m»\033[0m ",
		AutoComplete: commands.List.Build(),
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

		var inputErr = ParseInput(line)
		if inputErr != nil {
			fmt.Printf("%s\n", inputErr)
			continue
		}
	}
}

func ParseInput(line string) error {
	line = strings.TrimSpace(line)
	var parts = strings.Split(line, " ")
	var inputs, inputErr = commands.List.ParseInput(parts, nil)

	if inputErr != nil {
		return inputErr
	}

	inputs.Execute()

	return nil
}
