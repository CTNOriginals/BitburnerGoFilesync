package commands

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

var CommandList []*Definition = make([]*Definition, 0)
var promtSymbol = "\033[34m\033[1m»\033[0m "

func usage(writer io.Writer, options *readline.PrefixCompleter) {
	io.WriteString(writer, "commands:\n")
	io.WriteString(writer, options.Tree("    "))
}

func CommandWatcher() {
	// fmt.Printf("commands: %v\n", CommandList)

	var options = buildOptions()

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

	// idk yet what this does
	log.SetOutput(cli.Stderr())
	usage(cli.Stderr(), options)

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
			usage(cli.Stderr(), options)
			continue
		}

		for _, def := range CommandList {
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

	// var scanner = bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	var cmd = strings.TrimSpace(scanner.Text())
	//
	// 	var def = GetCommandByTrigger(cmd)
	//
	// 	if def == nil {
	// 		fmt.Printf("Unknown command: %s\n", cmd)
	// 		goto endscan
	// 	}
	//
	// 	def.Execution()
	//
	// endscan:
	// 	fmt.Print(promtSymbol)
	// }
}

func buildOptions() *readline.PrefixCompleter {
	// var options = readline.PrefixCompleter{
	// 	Children: make([]readline.PrefixCompleterInterface, 0),
	// }
	var options = make([]readline.PrefixCompleterInterface, len(CommandList))

	for i, def := range CommandList {
		options[i] = readline.PcItem(def.Triggers[0], def.Options)
	}

	return readline.NewPrefixCompleter(options...)
}

func GetCommandByTrigger(trigger string) *Definition {
	for _, def := range CommandList {
		if slices.Contains(def.Triggers, trigger) {
			return def
		}
	}

	return nil
}
