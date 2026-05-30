package clogger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type FState func() bool
type MLogLevelState = MLogLevel[FState]
type MLogLevelPrefix = MLogLevel[TPrefixMask]

type SClog struct {
	Name       string
	PrefixMask TPrefixMask

	DefaultState   FState
	LogLevelState  MLogLevelState
	LogLevelPrefix MLogLevelPrefix

	logger *log.Logger
}

func (this *SClog) initialize() {
	this.logger = log.New(log.Default().Writer(), "", 0)
}

func (this *SClog) checkInit() {
	if this.logger == nil {
		this.initialize()
	}
}

func (this *SClog) GetStackTrace(skip int) string {
	var trace = make([]byte, 1<<16)
	var traceSize = runtime.Stack(trace, false)
	var lines = strings.Split(string(trace[:traceSize]), "\n")

	// skip this GetStackTrace()
	skip += 1
	// each 2 lines is 1 trace
	skip *= 2
	// skip the top trace line that doesnt hold info
	skip += 1

	if len(lines) < skip {
		return string(trace[:traceSize])
	}

	return strings.Join(lines[skip:], "\n")
}

func (this *SClog) PrintStackTrace() {
	this.checkInit()
	this.logger.Print(this.GetStackTrace(1), "\n")
}

func (this SClog) checkState(level TLogLevel) bool {
	var stateChecker = this.DefaultState

	if levelState := this.LogLevelState.Get(level); levelState != nil {
		stateChecker = *levelState
	}

	if stateChecker == nil {
		return true
	}

	return stateChecker()
}

func (this *SClog) print(level TLogLevel, msg ...any) {
	this.checkInit()

	if !this.checkState(level) {
		return
	}

	var prefix = this.PrefixMask
	if levelPrefix := this.LogLevelPrefix.Get(level); levelPrefix != nil {
		prefix = *levelPrefix
	}

	var prefixString = prefix.GetPrefix(this.Name, level)

	msg = append([]any{prefixString}, msg...)

	// Print out the message
	this.logger.Print(msg...)

	if level.Has(LogFatal) {
		var tracePrefixMask = PrefixDate |
			PrefixTime |
			PrefixFile |
			PrefixCall |
			PrefixNoColor

		var tracePrefix = tracePrefixMask.GetPrefix(this.Name, level)
		this.logger.Printf("%s\n%s\n", tracePrefix, this.GetStackTrace(2))
	}
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
