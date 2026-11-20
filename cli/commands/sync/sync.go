package cmdsync

import (
	"os"
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli/commands"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "cmdsync",
})

var warning = "[WARNING] Force overwrites any existing files on bitburner."

var def = commands.Definition{
	Name: "sync",
	Description: []string{
		"Sync all your local files to bitburner at once.",
		warning,
	},

	Options: commands.TList{
		{
			Name: "confirm",
			Description: []string{
				"Whith this option, the sync will execute the operation.",
				warning,
			},
			ExpectValue: false,
		},
	},

	Execution: func(args commands.TInputList, self int) {
		clog.Debugf("%v\n%d", args, self)
		var pathMap = map[string]string{}

		utils.ForEachFileInDirRecursive(utils.GetAbsolutePath(""), func(file os.FileInfo, dir string) {
			var path = filepath.Join(dir, file.Name())
			pathMap[path] = utils.ToBitburnerPath(path)
		})

		var dry = true

		for _, arg := range args {
			if arg.Def.Name == "confirm" {
				dry = false
				break
			}
		}

		clog.Info("[WARNING] All the following files will be synced to bitburner:")
		for source, dest := range pathMap {
			clog.Messagef("%s: %s", source, dest)

			if dry {
				continue
			}

			websocket.Client.Socket.PushFile(websocket.Params_PushFile{
				Filename: dest,
				Content:  string(utils.GetFileContentByPath(source)),
				Server:   "home",
			}, nil)
		}

		if dry {
			clog.Info(
				"No action was taken.",
				"\nTo confirm the sync, you have to confirm this action.",
				"\n", warning,
				"\nrun: 'sync confirm' to execute the sync",
			)
			return
		}
	},
}

func init() {
	commands.List.Push(&def)
}
