package clogger

import (
	"fmt"
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogLevel utils.TBitMask

const (
	LogInfo ELogLevel = 1 << iota
	LogDebug
	LogWarn
	LogError
	LogFatal

	LogNone    = 0
	LogMessage = LogInfo | LogDebug
	LogAlert   = LogWarn | LogError | LogFatal
	LogAll     = (1 << iota) - 1
)

type SClog struct {
	Prefix     string
	PrefixTime bool
	PrefixDate bool

	LogLevel ELogLevel

	logger *log.Logger
}

func (this SClog) isInitialized() bool {
	return this.logger != nil
}

// Flag bits:
//
//	0b01: Time
//	0b10: Date
func (this SClog) getFlag() int {
	var flag = 0

	if this.PrefixDate {
		flag = 2
	}

	if this.PrefixTime {
		flag += 1
	}

	return flag
}

func (this *SClog) initialize() {
	this.logger = log.New(log.Default().Writer(), this.Prefix, this.getFlag())
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

func (this SClog) Print(msg ...any) {
	this.print(msg...)
}

func (this SClog) Printf(format string, args ...any) {
	this.printf(format, args...)
}
