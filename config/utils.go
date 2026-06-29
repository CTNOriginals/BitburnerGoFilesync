package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func ValidateFilePath(relative string, path string) (string, error) {
	var fileExists = func() error {
		var _, err = os.Stat(path)
		return err
	}

	if filepath.IsAbs(path) {
		return path, fileExists()
	}

	var baseDir = filepath.Dir(constants.ConfigFilePath)

	path = filepath.Clean(filepath.Join(baseDir, path))

	if filepath.IsAbs(path) {
		return path, fileExists()
	}

	var abspath, err = filepath.Abs(path)

	if err != nil {
		return path, err
	}

	path = abspath

	return path, fileExists()
}

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

	clog.Infof("Validated and set: %s\n", dir)

	Values.Directory = dir
}
