package watcher

import (
	"os"
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

var patternPaths config.TConfigFilrPatterns

func generatePatternPaths() {
	var patterns = config.Values.FilePatterns
	var dir = config.Values.Directory

	patternPaths.Include = make([]string, len(patterns.Include))
	patternPaths.Exclude = make([]string, len(patterns.Exclude))

	for i, inc := range patterns.Include {
		patternPaths.Include[i] = filepath.Join(dir, inc)
	}

	for i, exc := range patterns.Exclude {
		patternPaths.Exclude[i] = filepath.Join(dir, exc)
	}
}

func filePatternFilter(path string, info os.FileInfo) bool {

	return true
}
