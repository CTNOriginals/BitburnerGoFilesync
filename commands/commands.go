package commands

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
	"github.com/chzyer/readline"
)

var CommandList []*Definition = make([]*Definition, 0)
var promtSymbol = "\033[34m\033[1m»\033[0m "

func usage(writer io.Writer, options *readline.PrefixCompleter) {
	io.WriteString(writer, "commands:\n")
	io.WriteString(writer, options.Tree("    "))
}

func CommandWatcher() {
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

	log.SetOutput(cli.Stderr())

	fmt.Printf("options: %s\n", optToString(options, 1))
	// fmt.Printf("FileItem: %s\n", optToString(ReadLine_FileItem, 0))
	// fmt.Printf("File Completion (managers/): %v\n", ReadLine_FileItem.Callback("prototype managers/"))

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
}

func buildOptions() *readline.PrefixCompleter {
	var options = make([]readline.PrefixCompleterInterface, len(CommandList))

	for i, def := range CommandList {
		options[i] = readline.PcItem(def.Triggers[0], def.BuildOptions())
	}

	return readline.NewPrefixCompleter(options...)
}

func optToString(opt readline.PrefixCompleterInterface, depth int) string {
	var lines = []string{}

	var keys = ctnstruct.Keys(opt)
	var vals = ctnstruct.Values(opt)

	var callback string = fmt.Sprintf("%v", vals[slices.Index(keys, "Callback")])

	if callback != "<nil>" {
		return fmt.Sprintf("Callback: %v", callback)
	}

	var children = make([]string, len(opt.GetChildren()))
	for i, child := range opt.GetChildren() {
		var name = string(child.GetName())
		if name == "" {
			name = "Dynamic"
		}
		children[i] = fmt.Sprintf("%s: %s", name, optToString(child, depth+1))
	}

	if len(children) > 0 {
		lines = append(lines, fmt.Sprintf("[\n%s\n]", ctnstring.Indent(strings.Join(children, ",\n"), depth+1, " ")))
	}

	if len(lines) == 0 {
		return "{ }"
	} else if len(lines) == 1 {
		return lines[0]
	}

	var str = ctnstring.Indent(strings.Join(lines, "\n"), depth, " ")

	return fmt.Sprintf("{\n%s\n}", str)
}

func GetCommandByTrigger(trigger string) *Definition {
	for _, def := range CommandList {
		if slices.Contains(def.Triggers, trigger) {
			return def
		}
	}

	return nil
}
