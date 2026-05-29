package clogger

import (
	"log"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name: "test",
		PrefixMask: PrefixLevel |
			PrefixName |
			PrefixTime |
			PrefixDate |
			PrefixFile,

		DefaultState: func() bool {
			return constants.Debug
		},

		LogLevelState: MLogLevelState{
			LogInfo: func() bool {
				var now = time.Now()
				return now.Second()%2 == 0
			},
			LogError | LogFatal: func() bool {
				return true
			},
		},
	}

	clog.Info("hello world!")
	clog.Debugf("Levels: '%v'", clog.LogLevelState)
	clog.Error("something gone wrong!\nyou better fix it...")
	clog.Fatalf("fat alf done it again: %v", clog)
}
