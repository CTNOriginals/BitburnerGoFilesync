package clogger

import "log"

func TestLogger() {
	log.Printf("\n--- CLOGGER TEST ---\n")
	var clog = SClog{Prefix: "test: "}
	clog.Print("hello")
	clog.Printf("%s=%d", "x", 1)

	// fmt.Println(clog.getFlag())
	//
	// var l2 = SClog{PrefixTime: true, PrefixDate: true}
	// fmt.Println(l2.getFlag())
}
