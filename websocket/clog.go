package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

var clog = clogger.SClog{
	Name:       "websocket",
	PrefixMask: clogger.PrefixCall,

	LogLevelState: clogger.MLogLevelState{
		clogger.LogDebug: func() bool {
			return constants.Debug
		},
	},
}
