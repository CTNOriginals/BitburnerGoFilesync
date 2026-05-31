package clogger

import (
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogLevel = utils.TBitMask
type TLogLevel ELogLevel

const (
	LogInfo TLogLevel = 1 << iota
	// Special level: does not include a prefix
	LogMessage
	LogDebug
	LogError
	LogFatal
)

// The order/priority of each log level from high to low.
var LogLevelOrder = []TLogLevel{
	LogFatal,
	LogError,
	LogInfo,
	LogMessage,
	LogDebug,
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
