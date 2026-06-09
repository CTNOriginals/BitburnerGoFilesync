package watcher

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type SFilePatterns struct {
	include []string
	exclude []string

	// All the file paths that match Include and not Exclude
	//  map[dir]matches
	matches map[string][]string
}

func (this SFilePatterns) validatePatterns() error {
	for _, pattern := range append(this.include, this.exclude...) {
		var valid = doublestar.ValidatePathPattern(pattern)

		if !valid {
			return fmt.Errorf("Invalid file pattern: %s", pattern)
		}
	}

	return nil
}

func (this *SFilePatterns) setPatterns(inc []string, exc []string) error {
	// If include is empty, include anything
	if len(inc) == 0 {
		inc = append(inc, "**/*")
	}

	this.include = inc
	this.exclude = exc

	return this.validatePatterns()
}

func NewFilePatterns(inc []string, exc []string) (*SFilePatterns, error) {
	var fp = SFilePatterns{
		matches: make(map[string][]string),
	}

	var err = fp.setPatterns(inc, exc)

	return &fp, err
}

// Clears the existing matches
func (this *SFilePatterns) Refresh() {
	this.matches = make(map[string][]string)
}

// Sets the patterns again and refreshes the results
func (this *SFilePatterns) UpdatePatterns(inc []string, exc []string) error {
	this.Refresh()
	return this.setPatterns(inc, exc)
}

func (this SFilePatterns) getValidPaths(dir string) []string {
	var inc = make([]string, 0)
	var exc = make([]string, 0)

	for _, pattern := range this.exclude {
		var patternPath = filepath.Join(dir, pattern)
		var matches, _ = doublestar.FilepathGlob(patternPath)

		exc = append(exc, matches...)
	}

	for _, pattern := range this.include {
		var patternPath = filepath.Join(dir, pattern)
		var matches, _ = doublestar.FilepathGlob(patternPath)

		for _, match := range matches {
			if slices.Contains(exc, match) || slices.Contains(inc, match) {
				continue
			}

			inc = append(inc, match)
		}
	}

	return inc
}

func (this *SFilePatterns) GetValidPaths(dir string) []string {
	// NOTE: this may be nice for performance,
	// but currently too complex to implement
	var matches, exists = this.matches[dir]
	if exists {
		return matches
	}
	this.matches[dir] = this.getValidPaths(dir)
	return this.matches[dir]

	return this.getValidPaths(dir)
}

func (this *SFilePatterns) IsValidPath(dir string, path string) bool {
	var validPaths = this.GetValidPaths(dir)

	if !strings.HasPrefix(path, dir) {
		path = filepath.Join(dir, path)
	}

	if !slices.Contains(validPaths, path) {
		clog.Errorf("Invalid: %s", path)
	}

	return slices.Contains(validPaths, path)
}
