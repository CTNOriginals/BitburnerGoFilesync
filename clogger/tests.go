package clogger

import "log"

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name:       "test: ",
		PrefixMask: PrefixName | PrefixTime,
	}
	clog.Info("hello world!")
	clog.Debugf("Levels: '%v'", clog)
}
