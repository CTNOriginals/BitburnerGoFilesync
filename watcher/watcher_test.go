package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type TEventCall struct {
	event EFileEvent
	path  string
}

func TestEmitDispatchesCorrectHandler(t *testing.T) {
	var recorded string
	var handler = MEventHandler{
		OnFileCreate: func(path string) { recorded = path },
	}

	handler.Emit(OnFileCreate, "test.ts")

	if recorded != "test.ts" {
		t.Errorf("Emit recorded %q, want %q", recorded, "test.ts")
	}
}

func TestEmitPanicsOnMissingKey(t *testing.T) {
	var handler = MEventHandler{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing event handler")
		}
	}()

	handler.Emit(OnFileCreate, "test.ts")
}

func TestScanDetectsModification(t *testing.T) {
	var dir = t.TempDir()
	var filePath = filepath.Join(dir, "test.ts")
	os.WriteFile(filePath, []byte("original"), 0644)

	var past = time.Now().Add(-time.Hour)
	os.Chtimes(filePath, past, past)

	var savedState = fileStateMap
	var savedHandler = fileEventHandler
	t.Cleanup(func() {
		fileStateMap = savedState
		fileEventHandler = savedHandler
	})

	fileStateMap = MFileState{filePath: past}

	var recorded []TEventCall
	fileEventHandler = MEventHandler{
		OnFileModify: func(path string) { recorded = append(recorded, TEventCall{OnFileModify, path}) },
		OnFileDelete: func(path string) { recorded = append(recorded, TEventCall{OnFileDelete, path}) },
	}

	os.WriteFile(filePath, []byte("modified"), 0644)

	scan()

	if len(recorded) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(recorded), recorded)
	}
	if recorded[0].event != OnFileModify {
		t.Errorf("expected OnFileModify, got %v", recorded[0].event)
	}
	if recorded[0].path != filePath {
		t.Errorf("expected path %q, got %q", filePath, recorded[0].path)
	}

	if fileStateMap[filePath] == past {
		t.Error("scan() did not update fileStateMap with new modtime")
	}
}

func TestScanDetectsDeletion(t *testing.T) {
	var dir = t.TempDir()
	var filePath = filepath.Join(dir, "test.ts")
	os.WriteFile(filePath, []byte("content"), 0644)

	var savedState = fileStateMap
	var savedHandler = fileEventHandler
	t.Cleanup(func() {
		fileStateMap = savedState
		fileEventHandler = savedHandler
	})

	fileStateMap = MFileState{filePath: time.Now()}

	var recorded []TEventCall
	fileEventHandler = MEventHandler{
		OnFileDelete: func(path string) { recorded = append(recorded, TEventCall{OnFileDelete, path}) },
	}

	os.Remove(filePath)

	scan()

	if len(recorded) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(recorded), recorded)
	}
	if recorded[0].event != OnFileDelete {
		t.Errorf("expected OnFileDelete, got %v", recorded[0].event)
	}
	if recorded[0].path != filePath {
		t.Errorf("expected path %q, got %q", filePath, recorded[0].path)
	}

	if _, exists := fileStateMap[filePath]; exists {
		t.Error("scan() did not delete entry from fileStateMap")
	}
}

func TestScanSkipsUnchangedFiles(t *testing.T) {
	var dir = t.TempDir()
	var filePath = filepath.Join(dir, "test.ts")
	os.WriteFile(filePath, []byte("content"), 0644)

	var info, _ = os.Stat(filePath)
	var modTime = info.ModTime()

	var savedState = fileStateMap
	var savedHandler = fileEventHandler
	t.Cleanup(func() {
		fileStateMap = savedState
		fileEventHandler = savedHandler
	})

	fileStateMap = MFileState{filePath: modTime}

	var recorded []TEventCall
	fileEventHandler = MEventHandler{
		OnFileCreate: func(path string) { recorded = append(recorded, TEventCall{OnFileCreate, path}) },
		OnFileModify: func(path string) { recorded = append(recorded, TEventCall{OnFileModify, path}) },
		OnFileDelete: func(path string) { recorded = append(recorded, TEventCall{OnFileDelete, path}) },
	}

	scan()

	if len(recorded) != 0 {
		t.Errorf("expected 0 events for unchanged file, got %d: %v", len(recorded), recorded)
	}
}

func TestScanNoOpOnEmptyMap(t *testing.T) {
	var savedState = fileStateMap
	var savedHandler = fileEventHandler
	t.Cleanup(func() {
		fileStateMap = savedState
		fileEventHandler = savedHandler
	})

	fileStateMap = make(MFileState)

	var recorded []TEventCall
	fileEventHandler = MEventHandler{
		OnFileModify: func(path string) { recorded = append(recorded, TEventCall{OnFileModify, path}) },
		OnFileDelete: func(path string) { recorded = append(recorded, TEventCall{OnFileDelete, path}) },
	}

	scan()

	if len(recorded) != 0 {
		t.Errorf("expected 0 events for empty map, got %d", len(recorded))
	}
}
