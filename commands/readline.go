package commands

import (
	"os"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/chzyer/readline"
)

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := os.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var ReadLine_FileItem = readline.PcItemDynamic(listFiles(constants.WorkindDirectory))
