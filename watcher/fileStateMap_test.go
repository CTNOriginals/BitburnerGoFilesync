package watcher

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestPush(t *testing.T) {
	var dir = t.TempDir()
	var filePath = filepath.Join(dir, "test.ts")
	os.WriteFile(filePath, []byte("content"), 0644)

	var fsm = make(MFileState)
	var err = fsm.Push(filePath)

	if err != nil {
		t.Fatalf("Push error: %v", err)
	}

	var info, _ = os.Stat(filePath)
	if fsm[filePath] != info.ModTime() {
		t.Errorf("Push stored %v, want %v", fsm[filePath], info.ModTime())
	}
}

func TestPushNonExistentFile(t *testing.T) {
	var fsm = make(MFileState)
	var err = fsm.Push("/nonexistent/path.ts")

	if err == nil {
		t.Error("Push expected error for nonexistent file")
	}
}

func TestGetNewEntriesFindsMatchingFiles(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts"}, []string{})

	os.WriteFile(filepath.Join(dir, "a.ts"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.ts"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte("c"), 0644)

	var fsm = make(MFileState)
	var found []string
	fsm.GetNewEntries(dir, func(path string) {
		found = append(found, path)
	})

	if len(found) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(found), found)
	}

	sort.Strings(found)
	if filepath.Base(found[0]) != "a.ts" || filepath.Base(found[1]) != "b.ts" {
		t.Errorf("found files: %v, want [a.ts, b.ts]", found)
	}
}

func TestGetNewEntriesSkipsDirectories(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*"}, []string{})

	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "file.ts"), []byte("content"), 0644)

	var fsm = make(MFileState)
	var found []string
	fsm.GetNewEntries(dir, func(path string) {
		found = append(found, path)
	})

	if len(found) != 1 || filepath.Base(found[0]) != "file.ts" {
		t.Errorf("found files: %v, want [file.ts]", found)
	}
}

func TestGetNewEntriesRecursive(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"**/*.ts"}, []string{})

	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "nested.ts"), []byte("content"), 0644)

	var fsm = make(MFileState)
	var found []string
	fsm.GetNewEntries(dir, func(path string) {
		found = append(found, path)
	})

	if len(found) != 1 || filepath.Base(found[0]) != "nested.ts" {
		t.Errorf("found files: %v, want [nested.ts]", found)
	}
}

func TestGetNewEntriesSkipsExisting(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*.ts"}, []string{})

	var existingPath = filepath.Join(dir, "existing.ts")
	var newPath = filepath.Join(dir, "new.ts")
	os.WriteFile(existingPath, []byte("content"), 0644)
	os.WriteFile(newPath, []byte("content"), 0644)

	var fsm = make(MFileState)
	var info, _ = os.Stat(existingPath)
	fsm[existingPath] = info.ModTime()

	var found []string
	fsm.GetNewEntries(dir, func(path string) {
		found = append(found, path)
	})

	if len(found) != 1 || filepath.Base(found[0]) != "new.ts" {
		t.Errorf("found files: %v, want [new.ts]", found)
	}
}

func TestGetNewEntriesSkipsExcluded(t *testing.T) {
	var dir = t.TempDir()
	setupPatternTest(t, dir, []string{"*"}, []string{"*.d.ts"})

	os.WriteFile(filepath.Join(dir, "test.ts"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(dir, "test.d.ts"), []byte("content"), 0644)

	var fsm = make(MFileState)
	var found []string
	fsm.GetNewEntries(dir, func(path string) {
		found = append(found, path)
	})

	if len(found) != 1 || filepath.Base(found[0]) != "test.ts" {
		t.Errorf("found files: %v, want [test.ts]", found)
	}
}
