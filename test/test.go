package test

import (
	"fmt"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func DoTest() {
	fmt.Printf("%s: %v\n", "Port", constants.Port)
	fmt.Printf("%s: %v\n", "WorkindDirectory", constants.WorkindDirectory)
	fmt.Printf("%s: %v\n", "BitburnerRoot", constants.BitburnerRoot)
	fmt.Printf("%s: %v\n", "IncludeFileExt", constants.IncludeFileExt)
	fmt.Printf("%s: %v\n", "FileScanDelay", constants.FileScanDelay)
	fmt.Printf("%s: %v\n", "NoWatcher", constants.NoWatcher)
	fmt.Printf("%s: %v\n", "NoServer", constants.NoServer)
	fmt.Printf("%s: %v\n", "KeepAlive", constants.KeepAlive)
}
