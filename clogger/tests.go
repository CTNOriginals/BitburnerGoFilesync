package clogger

import (
	"log"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name:       "test",
		PrefixMask: PrefixLevel | PrefixName,

		DefaultState: func() bool {
			return constants.Debug
		},

		LogLevelState: MLogLevelState{
			LogInfo: func() bool {
				var now = time.Now()
				return now.Second()%2 == 0
			},
			LogError | LogFatal: func() bool {
				var now = time.Now()
				return now.Second()%2 != 0
			},
		},

		LogLevelPrefix: MLogLevelPrefix{
			LogDebug: PrefixFile | PrefixCall,
			LogError | LogFatal: PrefixLevel |
				PrefixName |
				PrefixTime |
				PrefixDate |
				PrefixFile,
		},
	}

	clog.Info("hello world!")
	clog.Debug("foo bar baz")
	clog.Error("something gone wrong!\nyou better fix it...")
	clog.Fatalf("fat alf done it again: %v", clog)
}
