package test

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func DoTest() {
	config.Initialize()
	println("\n")
	fmt.Printf("%s: %v\n", "Port", config.Values.Port)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	fmt.Printf("%s: %v\n", "BitburnerRoot", config.Values.Directory)
	fmt.Printf("%s: %v\n", "IncludeFileExt", config.Values.FilePatterns.Include)
	fmt.Printf("%s: %v\n", "FileScanDelay", config.Values.FilePatterns.Exclude)
	fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
	println("")

	// var content = utils.GetFileContentByPath("proto.ts")
	// fmt.Print(string(content))
}
