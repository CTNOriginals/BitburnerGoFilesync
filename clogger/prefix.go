package clogger

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/utils"
)

type EPrefixMask = utils.TBitMask
type TPrefixMask EPrefixMask

const (
	PrefixName TPrefixMask = 1 << iota
	PrefixLevel
	PrefixDate
	PrefixTime
	PrefixFile // TODO:
)

func (this TPrefixMask) Mask() EPrefixMask {
	return EPrefixMask(this)
}
