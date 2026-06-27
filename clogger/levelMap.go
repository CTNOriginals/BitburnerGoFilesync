package clogger

type MLogLevel[T any] map[TLogLevel]T

func (this MLogLevel[T]) Get(level TLogLevel) *T {
	for lvl, val := range this {
		if lvl.Has(level) {
			return &val
		}
	}

	return nil
}

var LogLevelNames = MLogLevel[string]{
	LogInfo:  "Info",
	LogDebug: "Debug",
	LogError: "Error",
	LogFatal: "Fatal",
}

var LogLevelColors = MLogLevel[string]{
	LogInfo:  "\x1b[36m",
	LogDebug: "\x1b[34m",
	LogError: "\x1b[31m",
	LogFatal: "\x1b[41;37;1m",
}
