package clogger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctntype "github.com/CTNOriginals/CTNGoUtils/v2/types"
)

type EPrefixMask = ctntype.TBitMask
type TPrefixMask EPrefixMask

const (
	PrefixName TPrefixMask = 1 << iota
	PrefixLevel
	PrefixDate
	PrefixTime
	PrefixFile

	// The function info that called the log
	PrefixCall

	// -- Prefix Options --

	// Prevents the prefix content from being colored
	PrefixNoColor

	// Makes PrefixFile contain the full file path
	PrefixFullPath
)

var PrefixOrder = []TPrefixMask{
	PrefixDate,
	PrefixTime,
	PrefixFile,
	PrefixCall,
	PrefixLevel,
	PrefixName,
}

func (this TPrefixMask) Mask() EPrefixMask {
	return EPrefixMask(this)
}

// Returns the string form of a single prefix
func (this TPrefixMask) getPrefixPart(part TPrefixMask) string {
	var now = time.Now()

	switch part {
	case PrefixDate:
		var date = now.Format(time.DateOnly)
		return date[2:]
	case PrefixTime:
		return now.Format(time.TimeOnly)
	case PrefixFile:
		var _, file, line, ok = runtime.Caller(4)

		if !ok {
			return "unknown:0"
		}

		if this.Mask().Has(PrefixFullPath.Mask()) {
			return fmt.Sprintf("%s:%d", file, line)
		}

		var parts = strings.Split(file, "/")
		file = parts[len(parts)-1]

		return fmt.Sprintf("%s:%d", file, line)
	case PrefixCall:
		var pc, _, _, ok = runtime.Caller(4)
		var info = runtime.FuncForPC(pc)

		if !ok || info == nil {
			return "func.unknown"
		}

		// remove the redundant program entry path
		var fnName = strings.TrimPrefix(info.Name(), fmt.Sprintf("%s/", constants.PackageEntryPath))

		return fnName
	}

	log.Printf("Unknown prefix mask: %b\n", part)
	return ""
}

func (this TPrefixMask) GetPrefix(name string, level TLogLevel) string {
	var builder strings.Builder

	for _, prefix := range PrefixOrder {
		if !this.Mask().Has(prefix.Mask()) {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteRune(' ')
		}

		switch prefix {
		case PrefixName:
			builder.WriteString(name)
			continue
		case PrefixLevel:
			builder.WriteString(level.String())
			continue
		}

		builder.WriteString(this.getPrefixPart(prefix))
	}

	if builder.Len() == 0 {
		return ""
	}

	var color, exists = LogLevelColors[level]

	if !exists ||
		config.Values.Logging.NoColor ||
		this.Mask().Has(PrefixNoColor.Mask()) {
		return fmt.Sprintf("%s", builder.String())
	}

	return fmt.Sprintf(
		"%s%s%s",
		color,
		builder.String(),
		"\x1b[0m",
	)
}
