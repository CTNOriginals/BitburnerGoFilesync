package watcher

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/bmatcuk/doublestar/v4"
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

	if len(patterns.Include) == 0 {
		patternPaths.Include = append(patternPaths.Include, filepath.Join(dir, "**/*"))
	}
}

func getPatternMatches() []string {
	var matches = config.TConfigFilrPatterns{
		Include: make([]string, 0),
		Exclude: make([]string, 0),
	}

	for _, exc := range patternPaths.Exclude {
		var list, _ = doublestar.FilepathGlob(exc)
		matches.Exclude = append(matches.Exclude, list...)
	}

	for _, inc := range patternPaths.Include {
		var list, _ = doublestar.FilepathGlob(inc)

		for _, match := range list {
			if !slices.Contains(matches.Exclude, match) && !slices.Contains(matches.Include, match) {
				matches.Include = append(matches.Include, match)
			}
		}
	}

	return matches.Include
}

func filePatternFilter(path string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}

	var matches = getPatternMatches()

	return slices.Contains(matches, path)
}
