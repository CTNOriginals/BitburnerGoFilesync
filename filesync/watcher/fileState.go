package watcher

import (
	"fmt"
	"os"

	ctnmap "github.com/CTNOriginals/CTNGoUtils/v2/map"
)

type FileInfo struct {
	Path string
	Info os.FileInfo
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

func (this FileInfo) String() string {
	// return ctnstruct.ToString(this)
	return fmt.Sprintf("Name: %s,\nModTime: %s,", this.Info.Name(), this.Info.ModTime())
}

type FileStateMap map[string]*FileInfo

func (this FileStateMap) String() string {
	return ctnmap.ToStringFunc(this, func(val *FileInfo) string {
		return val.String()
	})
}
