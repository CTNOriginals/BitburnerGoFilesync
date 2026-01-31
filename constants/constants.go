package constants

import (
	"fmt"
	"os"
)

// var Port = "8080"
var WorkindDirectory, _ = os.Getwd()

// var BitburnerRoot = WorkindDirectory
// var IncludeFileExt = []string{"js", "ts", "txt"}
var ConfigFilePath = fmt.Sprintf("%s/%s", WorkindDirectory, "config.toml")

// The time to sleep until the next file scan in miliseconds
// var FileScanDelay = 100

var Debug = false
var NoWatcher = false
var NoServer = false

var KeepAlive = false
