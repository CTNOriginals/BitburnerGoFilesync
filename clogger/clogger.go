package clogger

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
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
	this.logger = log.Default()
}

func (this *SClog) checkInit() {
	if this.logger == nil {
		this.initialize()
	}
}

// Makes a copy of this and applies all non nill
// fields in override to the copy and returns it.
//
// TODO: in any map field, apply the override keys that are present
// and leave the ones that are not defined per key in that map.
// Do make sure that no keys are added that conflict with existing ones.
func (this SClog) Clone(override SClog) SClog {
	if override.Name != "" {
		this.Name = override.Name
	}
	if override.PrefixMask != 0 {
		this.PrefixMask = override.PrefixMask
	}
	if override.DefaultState != nil {
		this.DefaultState = override.DefaultState
	}
	if override.LogLevelState != nil {
		this.LogLevelState = override.LogLevelState
	}
	if override.LogLevelPrefix != nil {
		this.LogLevelPrefix = override.LogLevelPrefix
	}

	this.logger = nil

	return this
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

	if level == LogMessage {
		this.logger.Print(msg...)
		return
	}

	var prefix = this.PrefixMask
	if levelPrefix := this.LogLevelPrefix.Get(level); levelPrefix != nil {
		prefix = *levelPrefix
	}

	var prefixString = prefix.GetPrefix(this.Name, level)
	if len(prefixString) > 0 {
		prefixString += ": "
	}

	msg = append([]any{prefixString}, msg...)

	var builder strings.Builder
	var message = fmt.Sprint(msg...)
	message = strings.TrimRight(message, "\n")
	var lines = strings.Split(message, "\n")

	builder.WriteString(lines[0])

	// Indent all extra lines after the initial by 2 spaces
	if len(lines) > 1 {
		builder.WriteRune('\n')
		var body = strings.Join(lines[1:], "\n")
		builder.WriteString(ctnstring.Indent(body, 2, " "))
	}

	this.logger.Print(builder.String())

	if level.Has(LogFatal) {
		var tracePrefixMask = PrefixDate |
			PrefixTime |
			PrefixFile |
			PrefixFullPath |
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

func (this *SClog) Message(msg ...any) {
	this.print(LogMessage, msg...)
}
func (this *SClog) Messagef(format string, args ...any) {
	this.print(LogMessage, fmt.Sprintf(format, args...))
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
