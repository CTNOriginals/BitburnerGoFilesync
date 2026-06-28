package constants

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var WorkindDirectory, _ = os.Getwd()

// The programs root package entry point
var PackageEntryPath = "github.com/CTNOriginals/BitburnerGoFilesync"

var ConfigFilePath = fmt.Sprintf("%s/%s", WorkindDirectory, "config.toml")

var Debug = false

var (
	NoWatcher = false
	NoServer  = false
	NoCli     = false
)

func init() {
	var pc, _, _, _ = runtime.Caller(0)
	var entryPath = runtime.FuncForPC(pc).Name()
	var parts = strings.Split(entryPath, "/")

	PackageEntryPath = strings.Join(parts[0:len(parts)-1], "/")
}
