package clogger

import (
	"fmt"
	"log"
)

type FState func() bool

type SClog struct {
	Name       string
	PrefixMask TPrefixMask

	DefaultState  FState
	LogLevelState MLogLevelState

	logger *log.Logger
}

func (this SClog) isInitialized() bool {
	return this.logger != nil
}

func (this *SClog) initialize() {
	this.logger = log.New(log.Default().Writer(), "", 0)
}

func (this *SClog) print(level TLogLevel, msg ...any) {
	if !this.isInitialized() {
		this.initialize()
	}

	var prefix = this.PrefixMask.GetPrefix(this.Name, level)
	msg = append([]any{prefix}, msg...)

	// Print out the message
	this.logger.Print(msg...)
}

// -- Log Level Functions --

func (this *SClog) Info(msg ...any) {
	this.print(LogInfo, msg...)
}
func (this *SClog) Infof(format string, args ...any) {
	this.print(LogInfo, fmt.Sprintf(format, args...))
}

func (this *SClog) Debug(msg ...any) {
	this.print(LogDebug, msg...)
}
func (this *SClog) Debugf(format string, args ...any) {
	this.print(LogDebug, fmt.Sprintf(format, args...))
}

func (this *SClog) Error(msg ...any) {
	this.print(LogError, msg...)
}
func (this *SClog) Errorf(format string, args ...any) {
	this.print(LogError, fmt.Sprintf(format, args...))
}

func (this *SClog) Fatal(msg ...any) {
	this.print(LogFatal, msg...)
}
func (this *SClog) Fatalf(format string, args ...any) {
	this.print(LogFatal, fmt.Sprintf(format, args...))
}
