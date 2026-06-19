package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

// Calls fn for each file in dir.
// Does not recurse into child directories.
// Also returns directories.
func ForEachFileInDir(dir string, fn func(file os.FileInfo)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		clog.Error(err)
		return
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			clog.Error(err)
			continue
		}

		fn(info)
	}
}

// Calls fn for each file in dir recursivly
func ForEachFileInDirRecursive(dir string, fn func(file os.FileInfo, dir string)) {
	ForEachFileInDir(dir, func(file os.FileInfo) {
		if file.IsDir() {
			childDir := fmt.Sprintf("%s/%s", dir, file.Name())
			ForEachFileInDirRecursive(childDir, fn)
			return
		}

		fn(file, dir)
	})
}

func GetRelativePath(path string) string {
	return strings.TrimPrefix(path, config.Values.Directory)
}

func GetAbsolutePath(path string) string {
	return filepath.Join(config.Values.Directory, path)
}

func ToBitburnerPath(path string) string {
	path = filepath.ToSlash(path)
	path = GetRelativePath(path)
	return path
}

// The path needs to be relative the the bitburner dir
func GetFileContentByPath(path string) []byte {
	if !ctnfile.FileExists(path) {
		clog.Errorf("File does not exist: %s\n", path)
		return []byte{}
	}

	return ctnfile.GetFileBytes(path)
}
