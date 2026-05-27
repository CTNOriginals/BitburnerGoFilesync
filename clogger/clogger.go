package clogger

import (
	"fmt"
	"log"
)

type SClog struct {
	Name   string
	Prefix ELogPrefix

	LogLevel TLogLevel
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

func (this *SClog) print(level TLogLevel, msg ...any) {
	if !this.isInitialized() {
		this.initialize()
	}

	var prefix = this.logger.Prefix()
	// Set the prefix to include the log level
	this.logger.SetPrefix(fmt.Sprintf("%s %s", level.String(), prefix))

	// Print out the message
	this.logger.Print(msg...)

	// Reset the Prefix to what it was before the level was added
	this.logger.SetPrefix(prefix)
}

func (this *SClog) printf(level TLogLevel, format string, args ...any) {
	this.print(level, fmt.Sprintf(format, args...))
}

func (this *SClog) Info(msg ...any) {
	this.print(LogInfo, msg...)
}
func (this *SClog) Infof(format string, args ...any) {
	this.printf(LogInfo, format, args...)
}
