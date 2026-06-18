package watcher

import (
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/bmatcuk/doublestar/v4"
)

var patternPaths config.TConfigFilrPatterns
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

	// clog.Debugf("Updated cache:\n%s", strings.Join(strings.Split(fmt.Sprintf("%v", pathCache), " /home/ctn/code/bitburner/testdir/"), "\n"))
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

	// NOTE: this may not be an error worthy case
	// as include and exclude can be set to not catch all patterns
	// so that all non-described file paths are concidered to be excluded
	// TODO: test this
	clog.Errorf("Path did not get added to filter cache after updating it. The path will be excluded: %s", path)

	pathCache[path] = false

	return false
}
