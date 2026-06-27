package clogger

import (
	"log"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")

	var clog = Default.Clone(SClog{
		Name: "clogger-test",
	})

	clog.Fatalf("logger obj: \n%s", ctnstruct.ToString(clog))
	clog.Debug()
	clog.Messagef("-- Message --\n%s", ctnstruct.ToString(clog))

	testAllFeatures()
}

func testAllFeatures() {
	var gclog = SClog{
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

	gclog.Info("hello world!")
	gclog.Debug("foo bar baz")
	gclog.Error("something gone wrong!\nyou better fix it...")
	gclog.Fatalf("fat alf done it again: %v", gclog)
}
