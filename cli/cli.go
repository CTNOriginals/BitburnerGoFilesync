package cli

import (
	"errors"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/chzyer/readline"
)

var clog = clogger.Default.Clone(clogger.SClog{})

func CommandWatcher() {
	var cli, err = readline.NewEx(&readline.Config{
		Prompt:              "\033[34m\033[1m»\033[0m ",
		AutoComplete:        commands.List.Build(),
		HistorySearchFold:   true,
		ForceUseInteractive: true,
	})

	if err != nil {
		clog.Fatalf("Commands readline error:\n%v", err)
	}

	var logWriter = log.Writer()
	log.SetOutput(cli.Stderr())

	defer func() {
		cli.Close()
		log.SetOutput(logWriter)
	}()

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
			clog.Errorf("%s", inputErr)
			continue
		}
	}
}

func GetInputSegments(input string) ([]string, error) {
	var segments = make([]string, 0)
	var segmentChars = "\"'`"
	var parts = strings.Split(input, " ")
	var char rune = 0
	var index = 0

	for i := 0; i < len(parts); i++ {
		var part = parts[i]
		var quotes = make([]rune, 0)

		for j := 0; j < len(part); j++ {
			var c = rune(part[j])

			if j > 0 && part[j-1] == '\\' {
				continue
			}

			if strings.ContainsRune(segmentChars, c) {
				quotes = append(quotes, c)
			}
		}

		var hasQuote = len(quotes) > 0

		switch {
		case char == 0 && !hasQuote:
			segments = append(segments, part)
			index += 1
		case char != 0:
			if slices.Contains(quotes, char) {
				char = 0
				segments = append(segments, strings.Join(parts[index:i+1], " "))
				index = i + 1
			}
		case hasQuote:
			char = quotes[0]
		}
	}

	if index != len(parts) {
		return nil, errors.New("Unable to parse input segments due to unmatched quote\n")
	}

	return segments, nil
}

func ParseInput(line string) error {
	line = strings.TrimSpace(line)
	var segments, err = GetInputSegments(line)

	if err != nil {
		return err
	}

	var inputs, inputErr = commands.List.ParseInput(segments[0], segments[1:]...)

	if inputErr != nil {
		return inputErr
	}

	inputs.Execute()

	return nil
}
