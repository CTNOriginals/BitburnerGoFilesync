package clogger

import "log"

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name: "test",
		PrefixMask: PrefixLevel |
			PrefixName |
			PrefixTime |
			PrefixDate |
			PrefixFile,
	}
	clog.Info("hello world!")
	clog.Debugf("Levels: '%v'", clog)
}
