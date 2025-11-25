package handler

import (
	"filesync/constants"
	"fmt"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/utils/file"
)

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
