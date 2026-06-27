package watcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func setupPatternTest(t *testing.T, dir string, include, exclude []string) {
	t.Helper()

	var saved = *config.Values
	t.Cleanup(func() { *config.Values = saved })

	config.Values.Directory = dir
	config.Values.FilePatterns = config.TConfigFilrPatterns{
		Include: include,
		Exclude: exclude,
	}

	generatePatternPaths()
}

func TestGeneratePatternPathsDefaultInclude(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{}, []string{})

	if len(patternPaths.Include) != 1 {
		t.Fatalf("expected 1 include pattern, got %d", len(patternPaths.Include))
	}

	var want = filepath.Join(dir, "**/*")
	if patternPaths.Include[0] != want {
		t.Errorf("include[0] = %q, want %q", patternPaths.Include[0], want)
	}
}

func TestGeneratePatternPathsJoinsDir(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts", "*.js"}, []string{})

	if len(patternPaths.Include) != 2 {
		t.Fatalf("expected 2 include patterns, got %d", len(patternPaths.Include))
	}

	var want0 = filepath.Join(dir, "*.ts")
	var want1 = filepath.Join(dir, "*.js")
	if patternPaths.Include[0] != want0 {
		t.Errorf("include[0] = %q, want %q", patternPaths.Include[0], want0)
	}
	if patternPaths.Include[1] != want1 {
		t.Errorf("include[1] = %q, want %q", patternPaths.Include[1], want1)
	}
}

func TestGeneratePatternPathsExcludes(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*"}, []string{"*.d.ts"})

	if len(patternPaths.Exclude) != 1 {
		t.Fatalf("expected 1 exclude pattern, got %d", len(patternPaths.Exclude))
	}

	var want = filepath.Join(dir, "*.d.ts")
	if patternPaths.Exclude[0] != want {
		t.Errorf("exclude[0] = %q, want %q", patternPaths.Exclude[0], want)
	}
}

func TestFilePatternFilterInclude(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts"}, []string{})

	var tsPath = filepath.Join(dir, "test.ts")
	os.WriteFile(tsPath, []byte("content"), 0644)

	if !filePatternFilter(tsPath) {
		t.Errorf("filePatternFilter(%q) = false, want true", tsPath)
	}
}

func TestFilePatternFilterExclude(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*"}, []string{"*.d.ts"})

	var tsPath = filepath.Join(dir, "test.ts")
	var dtsPath = filepath.Join(dir, "test.d.ts")
	os.WriteFile(tsPath, []byte("content"), 0644)
	os.WriteFile(dtsPath, []byte("content"), 0644)

	if !filePatternFilter(tsPath) {
		t.Errorf("filePatternFilter(%q) = false, want true", tsPath)
	}
	if filePatternFilter(dtsPath) {
		t.Errorf("filePatternFilter(%q) = true, want false", dtsPath)
	}
}

func TestFilePatternFilterUnknown(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts"}, []string{})

	var unknownPath = filepath.Join(dir, "unknown.txt")

	if filePatternFilter(unknownPath) {
		t.Errorf("filePatternFilter(%q) = true, want false", unknownPath)
	}
}

func TestFilePatternFilterCache(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts"}, []string{})

	var tsPath = filepath.Join(dir, "test.ts")
	os.WriteFile(tsPath, []byte("content"), 0644)

	if !filePatternFilter(tsPath) {
		t.Fatal("expected true on first call")
	}

	os.Remove(tsPath)

	if !filePatternFilter(tsPath) {
		t.Errorf("expected cached result true after file removal")
	}
}
