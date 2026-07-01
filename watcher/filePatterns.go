package watcher

import (
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/bmatcuk/doublestar/v4"
)

var patternPaths config.SConfigFilePatterns
var pathCache map[string]bool

func generatePatternPaths() {
	var patterns = config.Values.FilePatterns
	var dir = config.Values.Directory

	patternPaths.Include = make([]string, len(patterns.Include))
	patternPaths.Exclude = make([]string, len(patterns.Exclude))
	pathCache = make(map[string]bool)

	for i, inc := range patterns.Include {
		patternPaths.Include[i] = filepath.Join(dir, inc)
	}

	for i, exc := range patterns.Exclude {
		patternPaths.Exclude[i] = filepath.Join(dir, exc)
	}

	if len(patterns.Include) == 0 {
		patternPaths.Include = append(patternPaths.Include, filepath.Join(dir, "**/*"))
	}
}

func updatePathCache() {
	var cacheList = func(list []string, state bool) {
		for _, path := range list {
			pathCache[path] = state
		}
	}

	for _, inc := range patternPaths.Include {
		var list, _ = doublestar.FilepathGlob(inc)
		cacheList(list, true)
	}
	for _, exc := range patternPaths.Exclude {
		var list, _ = doublestar.FilepathGlob(exc)
		cacheList(list, false)
	}
}

func filePatternFilter(path string) bool {
	var cachedState, exists = pathCache[path]

	if exists {
		return cachedState
	}

	updatePathCache()

	cachedState, exists = pathCache[path]

	if exists {
		return cachedState
	}

	pathCache[path] = false

	return false
}
