package watcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
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
	file, err := os.Stat(utils.GetAbsolutePath(this.Path))

	if err != nil {
		clog.Errorf("watchet/GetInfo os.Stat: %v\n", err)
		return this.Info
	}

	return file
}

// Returns the path relative to the bitburner directory
func (this FileInfo) RelativePath() string {
	var split = strings.Split(this.Path, config.Values.Directory+"/")
	return split[len(split)-1]
}

type MFileState map[string]*FileInfo

func (this MFileState) String() string {
	return ctnmap.ToStringFunc(this, func(val *FileInfo) string {
		return val.String()
	})
}
