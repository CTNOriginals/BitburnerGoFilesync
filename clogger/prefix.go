package clogger

import "github.com/CTNOriginals/BitburnerGoFilesync/utils"

type ELogPrefix = utils.TBitMask

const (
	PrefixName ELogPrefix = 1 << iota
	PrefixLevel
	PrefixTime
	PrefixDate
)

func GetTimeDate(prefix ELogPrefix) utils.TBitMask {
	var flag utils.TBitMask = 0

	if prefix.Has(PrefixTime) {
		flag.Set(0b01)
	}
	if prefix.Has(PrefixDate) {
		flag.Set(0b10)
	}

	return flag
}
