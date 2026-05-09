package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/chzyer/readline"
)

func listFiles(filePath *string) func(string) []string {
	return func(line string) []string {
		var parts = strings.Split(line, " ")
		fmt.Printf("\n%s/%s\n", *filePath, parts[len(parts)-1])
		names := make([]string, 0)
		files, _ := os.ReadDir(*filePath)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var ReadLine_FileItem = readline.PcItemDynamic(listFiles(&config.Values.Directory))
