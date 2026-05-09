package commands

import (
	"os"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/chzyer/readline"
)

func Readline_FileComplete(filePath *string) func(string) []string {
	return func(line string) []string {
		// var parts = strings.Split(line, " ")
		var names = make([]string, 0)
		var files, _ = os.ReadDir(*filePath)

		// fmt.Printf("\n%s/%s\n", *filePath, parts[len(parts)-1])

		for _, file := range files {
			names = append(names, file.Name())
		}

		return names
	}
}

var ReadLine_FileItem = readline.PcItemDynamic(Readline_FileComplete(&config.Values.Directory))
