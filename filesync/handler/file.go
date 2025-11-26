package handler

import (
	"filesync/constants"
	"fmt"
	"os"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/utils/file"
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
			ForEachFileInDirRecursive(dir+file.Name(), fn)
			return
		}

		fn(file, dir)
	})
}

// The path needs to be relative the the bitburner dir
func GetFileByPath(path string) []rune {
	filePath := ctnfile.ParseFilePath(constants.BitburnerRelativePath + path)

	// println(filePath.String())

	if !ctnfile.FileExists(filePath.Full) {
		fmt.Printf("File does not exist: %s\n", path)
		return []rune{}
	}

	return ctnfile.GetFileRunes(filePath.Full)
}

// Will return the content with any string termenating character escapes.
//
// The output of this will be able to be passed in as a json value
// without it escaping out of its own falue field
func SanitizeFileContent(content []rune) []rune {
	sanitized := []rune{}

	for i, char := range content {
		switch char {
		case '\\':
			switch content[i+1] {
			case 'n', 't', 'r':
				//? Add an extra '\' to preserve the escaped string
				sanitized = append(sanitized, '\\', char)
			}
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
