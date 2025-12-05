package constants

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

var Port = "8080"
var WorkindDirectory, _ = os.Getwd()
var BitburnerRoot = WorkindDirectory
var IncludeFileExt = []string{"js", "ts", "txt"}

// The time to sleep until the next file scan in miliseconds
var FileScanDelay = 100

var NoWatcher = false
var NoServer = false
var KeepAlive = false

func SetBitburnerDir(dir string) {
	var isAbsolute = path.IsAbs(dir)

	if runtime.GOOS == "windows" {
		// checks if the second and third char of the path are ":/" or ":\"
		isAbsolute = dir[1] == ':' && (dir[2] == '/' || dir[2] == '\\')
	}

	if !isAbsolute {
		dir = fmt.Sprintf("%s/%s", WorkindDirectory, dir)
	}

	// Replace all back slashes (\) with forward ones (/)
	dir = strings.Join(strings.Split(dir, "\\"), "/")

	if dir[len(dir)-1] == '/' {
		dir = dir[0 : len(dir)-1]
	}

	dir = path.Clean(dir)

	fmt.Printf("Set the bitburner working directory to: %s\n", dir)

	BitburnerRoot = dir
}
