package watcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"

	ctnmap "github.com/CTNOriginals/CTNGoUtils/v2/map"
)

type FileInfo struct {
	Path string
	Info os.FileInfo
}

func (this FileInfo) String() string {
	// return ctnstruct.ToString(this)
	return fmt.Sprintf("Name: %s,\nModTime: %s,", this.Info.Name(), this.Info.ModTime())
}

// Gets the current os.FileInfo, not the info stored in this.info
func (this FileInfo) GetInfo() os.FileInfo {
	file, err := os.Stat(this.Path)

	if err != nil {
		println(err)
		return this.Info
	}

	return file
}

// Returns the path relative to the bitburner directory
func (this FileInfo) RelativePath() string {
	var split = strings.Split(this.Path, constants.BitburnerRoot+"/")
	return split[len(split)-1]
}

type MFileState map[string]*FileInfo

func (this MFileState) String() string {
	return ctnmap.ToStringFunc(this, func(val *FileInfo) string {
		return val.String()
	})
}
