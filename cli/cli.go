package cli

import (
	"io"
	"log"
	"slices"
	"strings"
	"unicode"

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

var quotes = []rune("\"'`")

func getInputQuoteRegions(input string) [][]int {
	var regions = make([][]int, 0)

	for i, char := range input {
		if !slices.Contains(quotes, char) || (i > 0 && input[i-1] == '\\') {
			continue
		}

		if len(regions) == 0 || regions[len(regions)-1][1] != -1 {
			regions = append(regions, []int{i, -1})
			continue
		}

		var current = regions[len(regions)-1]

		if rune(input[current[0]]) == char {
			current[1] = i
		}
	}

	// Clear any unmet quote regions
	for len(regions) > 0 && regions[len(regions)-1][1] == -1 {
		regions = regions[:len(regions)-1]
	}

	return regions
}

func GetInputSegments(input string) []string {
	var regions = getInputQuoteRegions(input)

	if len(regions) == 0 {
		return strings.Split(input, " ")
	}

	var segments = make([]string, 0)
	var orig = 0
	var reg = 0

	for i := 0; i < len(input); i++ {
		var char = rune(input[i])

		if reg <= len(regions)-1 && i == regions[reg][0] {
			i = regions[reg][1]
			reg += 1
			continue
		}

		if !unicode.IsSpace(char) {
			continue
		}

		segments = append(segments, input[orig:i])
		orig = i + 1
	}

	segments = append(segments, input[orig:])

	return segments
}

func ParseInput(line string) error {
	line = strings.TrimSpace(line)
	var segments = GetInputSegments(line)

	var inputs, inputErr = commands.List.ParseInput(segments[0], segments[1:]...)

	if inputErr != nil {
		return inputErr
	}

	inputs.Execute()

	return nil
}
