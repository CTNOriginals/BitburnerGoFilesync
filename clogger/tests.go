package clogger

import "log"

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name:     "test: ",
		Prefix:   PrefixName | PrefixTime,
		LogLevel: LogInfo | LogError,
	}
	clog.Info("hello world!")
	clog.Debugf("Levels: '%s'", clog.LogLevel.String())
}
