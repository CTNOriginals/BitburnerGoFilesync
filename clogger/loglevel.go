package clogger

import "github.com/CTNOriginals/BitburnerGoFilesync/utils"

type ELogLevel utils.TBitMask

const (
	LogInfo ELogLevel = 1 << iota
	LogDebug
	LogWarn
	LogError
	LogFatal

	LogNone    = 0
	LogMessage = LogInfo | LogDebug
	LogAlert   = LogWarn | LogError | LogFatal
	LogAll     = (1 << iota) - 1
)

type MLogLevelState map[ELogLevel]func() bool
