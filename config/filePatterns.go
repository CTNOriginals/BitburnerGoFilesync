package config

import (
	"fmt"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type SConfigFilePatterns struct {
	Include []string
	Exclude []string
}

func (this SConfigFilePatterns) ValidateValues() error {
	var invalidPatterns = make([]string, 0)

	for _, pattern := range append(Values.FilePatterns.Include, Values.FilePatterns.Exclude...) {
		if !doublestar.ValidatePathPattern(pattern) {
			invalidPatterns = append(invalidPatterns, pattern)
		}
	}

	if len(invalidPatterns) > 0 {
		return fmt.Errorf("The following file pattens are invalid:\n%s", strings.Join(invalidPatterns, "\n"))
	}

	return nil
}
