package clogger

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogPrefix = utils.TBitMask

const (
	PrefixName ELogPrefix = 1 << iota
	PrefixLevel
	PrefixDate
	PrefixTime
	// PrefixShortFile // TODO:
	// PrefixLongFile // TODO:

	// Put the name of the logger at the beginning of the line
	// instead of at the beginning of the message.
	PrefixNameAtStart
)

func PrefixGetLogFlag(prefix ELogPrefix) utils.TBitMask {
	var flag utils.TBitMask = 0

	flag.SetIf(log.Ldate, prefix.Has(PrefixDate))
	flag.SetIf(log.Ltime, prefix.Has(PrefixTime))
	flag.SetIf(log.Lmsgprefix, !prefix.Has(PrefixNameAtStart))

	return flag
}
