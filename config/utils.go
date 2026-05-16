package config

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func ValidateBitburnerDirectory(dir string) {
	var isAbsolute = path.IsAbs(dir)

	if runtime.GOOS == "windows" {
		// checks if the second and third char of the path are ":/" or ":\"
		isAbsolute = dir[1] == ':' && (dir[2] == '/' || dir[2] == '\\')
	}

	if !isAbsolute {
		dir = fmt.Sprintf("%s/%s", constants.WorkindDirectory, dir)
	}

	// Replace all back slashes (\) with forward ones (/)
	dir = strings.Join(strings.Split(dir, "\\"), "/")

	if dir[len(dir)-1] == '/' {
		dir = dir[0 : len(dir)-1]
	}

	dir = path.Clean(dir)

	log.Printf("Set the bitburner working directory to: %s\n", dir)

	Values.Directory = dir
}
