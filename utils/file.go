package utils

import (
	"fmt"
	"os"

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

func GetAbsolutePath(path string) string {
	return fmt.Sprintf("%s%s%s", config.Values.Directory, string(os.PathSeparator), path)
}

// The path needs to be relative the the bitburner dir
func GetFileContentByPath(path string) []byte {
	var filePath = GetAbsolutePath(path)

	if !ctnfile.FileExists(filePath) {
		clog.Errorf("File does not exist: %s\n", filePath)
		return []byte{}
	}

	return ctnfile.GetFileBytes(filePath)
}
