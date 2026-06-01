package utils

import "github.com/CTNOriginals/BitburnerGoFilesync/clogger"

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "utils",
})
