package clogger

import "log"

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{
		Name:     "test: ",
		Prefix:   PrefixName | PrefixDate,
		LogLevel: LogInfo | LogError,
	}
	clog.Info("hello world!")
	clog.Infof("Levels: '%s'", clog.LogLevel.String())

	// fmt.Println(clog.getFlag())
	//
	// var l2 = SClog{PrefixTime: true, PrefixDate: true}
	// fmt.Println(l2.getFlag())
}
