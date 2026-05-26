package clogger

import (
	"fmt"
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
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
	var flag utils.TBitMask = 0
	if this.Prefix.Has(PrefixTime) {
		flag.Set(0b01)
	}
	if this.Prefix.Has(PrefixDate) {
		flag.Set(0b10)
	}

	this.logger = log.New(log.Default().Writer(), this.Name, int(flag))
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
