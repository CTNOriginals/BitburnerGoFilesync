package utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

// Calls fn for each file in dir.
// Does not recurse into child directories.
// Also returns directories.
func ForEachFileInDir(dir string, fn func(file os.FileInfo)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		println(err)
		return
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			println(err)
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

// The path needs to be relative the the bitburner dir
func GetFileContentByPath(path string) []byte {
	var filePath = fmt.Sprintf("%s/%s", config.Values.Directory, path)

	if !ctnfile.FileExists(filePath) {
		fmt.Printf("utils/GetFileContentByPath: File does not exist: %s\n", filePath)
		return []byte{}
	}

	return ctnfile.GetFileBytes(filePath)
}

// Will return the content with any string termenating character escapes.
//
// The output of this will be able to be passed in as a json value
// without it escaping out of its own falue field
func SanitizeFileContent(content []byte) []byte {
	sanitized := []byte{}

	for _, char := range content {
		switch char {
		case '\\':
			sanitized = append(sanitized, '\\', char)

			// switch content[i+1] {
			// case 'n', 't', 'r', '/':
			// 	//? Add an extra '\' to preserve the escaped string
			// 	sanitized = append(sanitized, '\\', char)
			// }
		case '\r':
			continue
		case '"':
			sanitized = append(sanitized, '\\', char)
		case '\n':
			sanitized = append(sanitized, '\\', 'n')
		case '\t':
			sanitized = append(sanitized, '\\', 't')
		default:
			sanitized = append(sanitized, char)
		}
	}

	return sanitized
}

func SetBitburnerDir(dir string) {
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

	fmt.Printf("Set the bitburner working directory to: %s\n", dir)

	config.Values.Directory = dir
}
