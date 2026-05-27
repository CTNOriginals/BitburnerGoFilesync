package clogger

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogLevel = utils.TBitMask
type TLogLevel ELogLevel
type MLogLevel map[TLogLevel]string

const (
	LogInfo TLogLevel = 1 << iota
	LogDebug
	LogWarn
	LogError
	LogFatal

	LogNone    = 0
	LogMessage = LogInfo | LogDebug
	LogAlert   = LogWarn | LogError | LogFatal
	LogAll     = (1 << iota) - 1
)

var LogLevelTable = MLogLevel{
	LogInfo:  "Info",
	LogDebug: "Debug",
	LogWarn:  "Warn",
	LogError: "Error",
	LogFatal: "Fatal",
}

func (this TLogLevel) Mask() ELogLevel {
	return ELogLevel(this)
}

func (this TLogLevel) String() string {
	var str strings.Builder

	for typ, name := range LogLevelTable {
		if !this.Mask().Has(typ.Mask()) {
			continue
		}

		if str.Len() > 0 {
			str.WriteString(" ")
		}

		str.WriteString(name)
	}

	return str.String()
}

type MLogLevelState map[TLogLevel]func() bool
