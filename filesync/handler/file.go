package handler

import (
	"fmt"

	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/utils/file"
)

// The path needs to be relative the the bitburner dir
func GetFileByPath(path string) string {
	filePath := ctnfile.ParseFilePath("../bitburner/" + path)

	println(filePath.String())

	if !ctnfile.FileExists(filePath.Full) {
		fmt.Printf("File does not exist: %s\n", path)
		return ""
	}

	return string(ctnfile.GetFileRunes(filePath.Full))
}
