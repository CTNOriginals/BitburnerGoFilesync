package commands

import "github.com/CTNOriginals/BitburnerGoFilesync/clogger"

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "command",
})

var List = make(TList, 0)
