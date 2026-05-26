package clogger

import (
	"fmt"
	"log"
)

type SClog struct {
	Name   string
	Prefix ELogPrefix

	LogLevel ELogLevel
	LogState MLogLevelState

	logger *log.Logger
}

func (this SClog) isInitialized() bool {
	return this.logger != nil
}

func (this *SClog) initialize() {
	var name = ""

	if this.Prefix.Has(PrefixName) {
		name = this.Name
	}

	this.logger = log.New(log.Default().Writer(), name, int(PrefixGetLogFlag(this.Prefix)))
}

func (this *SClog) print(msg ...any) {
	if !this.isInitialized() {
		this.initialize()
	}

	this.logger.Print(msg...)
}

func (this *SClog) printf(format string, args ...any) {
	this.print(fmt.Sprintf(format, args...))
}

func (this *SClog) Print(msg ...any) {
	this.print(msg...)
}

func (this *SClog) Printf(format string, args ...any) {
	this.printf(format, args...)
}
