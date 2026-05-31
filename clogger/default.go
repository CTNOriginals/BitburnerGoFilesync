package clogger

import "github.com/CTNOriginals/BitburnerGoFilesync/constants"

var Default = SClog{
	Name:       "main",
	PrefixMask: PrefixCall,

	LogLevelState: MLogLevelState{
		LogDebug: func() bool {
			return constants.Debug
		},
	},

	LogLevelPrefix: MLogLevelPrefix{
		LogError | LogFatal: PrefixFile | PrefixCall,
	},
}
