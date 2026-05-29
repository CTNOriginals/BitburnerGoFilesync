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
	LogError
	LogFatal
)

var LogLevelTable = MLogLevel{
	LogInfo:  "Info",
	LogDebug: "Debug",
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

type MLogLevelState map[TLogLevel]FState
