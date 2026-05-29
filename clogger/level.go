package clogger

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogLevel = utils.TBitMask
type TLogLevel ELogLevel

const (
	LogInfo TLogLevel = 1 << iota
	LogDebug
	LogError
	LogFatal
)

// The order/priority of each log level from high to low.
var LogLevelOrder = []TLogLevel{
	LogFatal,
	LogError,
	LogInfo,
	LogDebug,
}

var LogLevelNames = map[TLogLevel]string{
	LogInfo:  "Info",
	LogDebug: "Debug",
	LogError: "Error",
	LogFatal: "Fatal",
}

func (this TLogLevel) Mask() ELogLevel {
	return ELogLevel(this)
}

func (this TLogLevel) Has(flag TLogLevel) bool {
	return this.Mask().Has(flag.Mask())
}

func (this TLogLevel) String() string {
	var str strings.Builder

	for _, typ := range LogLevelOrder {
		if !this.Mask().Has(typ.Mask()) {
			continue
		}

		if str.Len() > 0 {
			str.WriteString(" ")
		}

		str.WriteString(LogLevelNames[typ])
	}

	return str.String()
}

type MLogLevel[T any] map[TLogLevel]T

func (this MLogLevel[T]) Get(level TLogLevel) *T {
	for lvl, val := range this {
		if lvl.Has(level) {
			return &val
		}
	}

	return nil
}

type MLogLevelState = MLogLevel[FState]
type MLogLevelPrefix = MLogLevel[TPrefixMask]
