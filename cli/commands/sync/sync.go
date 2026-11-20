package cmdsync

import (
	"os"
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "cmdsync",
})

var def = commands.Definition{
	Name: "sync",
	Description: []string{
		"Sync all your local files to bitburner at once.",
	},

	Options: commands.TList{},

	Execution: func(args commands.TInputList, self int) {
		// clog.Debugf("Syncing: [%v] (%d)", args, self)
		utils.ForEachFileInDirRecursive(utils.GetAbsolutePath(""), func(file os.FileInfo, dir string) {
			var path = filepath.Join(dir, file.Name())
			clog.Debugf("Sync %s\nto %s", path, utils.ToBitburnerPath(path))
		})
	},
}

func init() {
	commands.List.Push(&def)
}
