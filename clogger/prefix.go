package clogger

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type ELogPrefix = utils.TBitMask
type TLogPrefix ELogPrefix

const (
	PrefixName TLogPrefix = 1 << iota
	PrefixLevel
	PrefixDate
	PrefixTime
	// PrefixShortFile // TODO:
	// PrefixLongFile // TODO:

	// Put the name of the logger at the beginning of the line
	// instead of at the beginning of the message.
	PrefixNameAtStart
)

func (this TLogPrefix) Mask() ELogPrefix {
	return ELogPrefix(this)
}

func (this TLogPrefix) GetLogFlag() utils.TBitMask {
	var flag utils.TBitMask = 0
	var mask = this.Mask()

	flag.SetIf(log.Ldate, mask.Has(PrefixDate.Mask()))
	flag.SetIf(log.Ltime, mask.Has(PrefixTime.Mask()))
	flag.SetIf(log.Lmsgprefix, !mask.Has(PrefixNameAtStart.Mask()))

	return flag
}
