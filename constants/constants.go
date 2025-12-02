package constants

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var WorkindDirectory, _ = os.Getwd()
var BitburnerRoot = WorkindDirectory

var NoWatcher = false
var NoServer = false
var KeepAlive = false

var IncludeFileExt = []string{"js", "ts", "txt"}

func SetBitburnerDir(dir string) {
	//! This may not work on linux
	if !(dir[1] == ':' && (dir[2] == '/' || dir[2] == '\\')) {
		dir = fmt.Sprintf("%s/%s", WorkindDirectory, dir)
	}

	// Replace all back slashes (\) with forward ones (/)
	dir = strings.Join(strings.Split(dir, "\\"), "/")

	if dir[len(dir)-1] == '/' {
		dir = dir[0 : len(dir)-1]
	}

	dir = path.Clean(dir)

	fmt.Printf("Set the bitburner working directory to: %s\n", dir)
}

const (
	Port = 8080
)
