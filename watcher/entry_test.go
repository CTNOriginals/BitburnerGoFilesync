package watcher

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	mustWriteFile(t, filepath.Join(dir, "file.txt"), []byte("hello"), 0644)

	mustMkdirAll(t, filepath.Join(dir, "subdir", "empty"), 0755)
	mustWriteFile(t, filepath.Join(dir, "subdir", "a.txt"), []byte("nested"), 0644)

	mustWriteFile(t, filepath.Join(dir, ".hidden"), []byte("dotfile"), 0644)

	noperm := filepath.Join(dir, "noperm")
	mustMkdirAll(t, noperm, 0755)
	mustWriteFile(t, filepath.Join(noperm, "secret.txt"), []byte("secret"), 0644)

	mustWriteFile(t, filepath.Join(dir, "notdir_test"), []byte("i am a file not a dir"), 0644)

	if runtime.GOOS != "windows" {
		chain := filepath.Join(dir, "chain")
		mustMkdirAll(t, chain, 0755)
		os.Symlink(filepath.Join(chain, "link2"), filepath.Join(chain, "link1"))
		os.Symlink(filepath.Join(chain, "link1"), filepath.Join(chain, "link2"))
	}

	return dir
}

func mustWriteFile(t *testing.T, path string, data []byte, perm os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, data, perm); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func mustMkdirAll(t *testing.T, path string, perm os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(path, perm); err != nil {
		t.Fatalf("MkdirAll(%s): %v", path, err)
	}
}

func mustChmod(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("Chmod(%s, %o): %v", path, mode, err)
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	fn()
}

// --- newEntry ---

func TestNewNode_File(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	if entry.info == nil {
		t.Fatal("expected non-nil info for existing file")
	}
	if entry.infoError != nil {
		t.Fatalf("expected nil error, got %v", entry.infoError)
	}
	if entry.info.Name() != "file.txt" {
		t.Errorf("expected name 'file.txt', got %q", entry.info.Name())
	}
}

func TestNewNode_Directory(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))

	if entry.info == nil {
		t.Fatal("expected non-nil info for existing directory")
	}
	if entry.infoError != nil {
		t.Fatalf("expected nil error, got %v", entry.infoError)
	}
	if !entry.info.IsDir() {
		t.Error("expected IsDir to be true")
	}
}

func TestNewNode_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	if entry.info != nil {
		t.Error("expected nil info for non-existent path")
	}
	if entry.infoError == nil {
		t.Fatal("expected error for non-existent path")
	}
	if _, ok := entry.infoError.(*os.PathError); !ok {
		t.Errorf("expected *os.PathError, got %T", entry.infoError)
	}
	if !os.IsNotExist(entry.infoError) {
		t.Errorf("expected not-exist error, got %v", entry.infoError)
	}
}

func TestNewNode_NotADirectory(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "notdir_test", "child")
	entry := newEntry(path)

	if entry.info != nil {
		t.Error("expected nil info for ENOTDIR path")
	}
	if entry.infoError == nil {
		t.Fatal("expected error for ENOTDIR path")
	}
}

func TestNewNode_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	noperm := filepath.Join(dir, "noperm")
	t.Cleanup(func() { os.Chmod(noperm, 0755) })
	mustChmod(t, noperm, 0000)

	path := filepath.Join(noperm, "secret.txt")
	entry := newEntry(path)

	if entry.info != nil {
		t.Error("expected nil info for permission-denied path")
	}
	if entry.infoError == nil {
		t.Fatal("expected error for permission-denied path")
	}
	if pathErr, ok := entry.infoError.(*os.PathError); ok {
		if !os.IsPermission(pathErr) {
			t.Errorf("expected permission error, got %v", pathErr.Err)
		}
	} else {
		t.Errorf("expected *os.PathError, got %T", entry.infoError)
	}
}

func TestNewNode_SymlinkLoop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink loop test requires unix")
	}
	dir := setupTestDir(t)
	path := filepath.Join(dir, "chain", "link1")
	entry := newEntry(path)

	if entry.info != nil {
		t.Error("expected nil info for symlink loop path")
	}
	if entry.infoError == nil {
		t.Fatal("expected error for symlink loop path")
	}
}

func TestNewNode_PathTooLong(t *testing.T) {
	dir := setupTestDir(t)
	longPath := dir
	for len(longPath) < 4096 {
		longPath = filepath.Join(longPath, "a")
	}
	entry := newEntry(longPath)

	if entry.info != nil {
		t.Error("expected nil info for path-too-long")
	}
	if entry.infoError == nil {
		t.Fatal("expected error for path-too-long")
	}
}

// --- Exists ---

func TestExists_Existing(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	if !entry.Exists() {
		t.Error("expected Exists() to be true for existing file")
	}
}

func TestExists_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	if entry.Exists() {
		t.Error("expected Exists() to be false for non-existent path")
	}
}

// --- IsDirectory ---

func TestIsDirectory_File(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	if entry.IsDirectory() {
		t.Error("expected IsDirectory to be false for a file")
	}
}

func TestIsDirectory_Directory(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))

	if !entry.IsDirectory() {
		t.Error("expected IsDirectory to be true for a directory")
	}
}

func TestIsDirectory_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { entry.IsDirectory() })
}

// --- GetPath ---

func TestGetPath(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	expected := filepath.Join(dir, "file.txt")
	if got := entry.GetPath(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGetPath_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { entry.GetPath() })
}

// --- IsModified ---

func TestIsModified_Unchanged(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	if entry.IsModified() {
		t.Error("expected IsModified to be false for unchanged file")
	}
}

func TestIsModified_Changed(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	entry := newEntry(path)

	time.Sleep(10 * time.Millisecond)

	mustWriteFile(t, path, []byte("modified"), 0644)

	if !entry.IsModified() {
		t.Error("expected IsModified to be true after file modification")
	}
}

func TestIsModified_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { entry.IsModified() })
}

// --- Update ---

func TestUpdate_NoChange(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	entry.Update()

	if entry.infoError != nil {
		t.Errorf("expected nil infoError after Update on unchanged file, got %v", entry.infoError)
	}
}

func TestUpdate_AfterModify(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	entry := newEntry(path)

	origModTime := entry.info.ModTime()

	time.Sleep(10 * time.Millisecond)
	mustWriteFile(t, path, []byte("modified"), 0644)

	entry.Update()

	if entry.info.ModTime() == origModTime {
		t.Error("expected ModTime to change after Update")
	}
	if entry.infoError != nil {
		t.Errorf("expected nil infoError, got %v", entry.infoError)
	}
}

func TestUpdate_AfterDelete(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	entry := newEntry(path)

	origName := entry.info.Name()

	os.Remove(path)

	entry.Update()

	if entry.infoError == nil {
		t.Error("expected infoError after file deletion")
	}
	if !os.IsNotExist(entry.infoError) {
		t.Errorf("expected not-exist error, got %v", entry.infoError)
	}
	if entry.info == nil {
		t.Fatal("expected info to be preserved after Update with error")
	}
	if entry.info.Name() != origName {
		t.Errorf("expected name %q to be preserved, got %q", origName, entry.info.Name())
	}
}

func TestUpdate_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	noperm := filepath.Join(dir, "noperm")
	path := filepath.Join(noperm, "secret.txt")
	entry := newEntry(path)

	origName := entry.info.Name()

	t.Cleanup(func() { os.Chmod(noperm, 0755) })
	mustChmod(t, noperm, 0000)

	entry.Update()

	if entry.infoError == nil {
		t.Error("expected infoError after permission denied")
	}
	if entry.info == nil {
		t.Fatal("expected info to be preserved")
	}
	if entry.info.Name() != origName {
		t.Errorf("expected name %q to be preserved, got %q", origName, entry.info.Name())
	}
}

// --- SortChildren ---

func TestSortChildren(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(dir)
	entry.children = []*SEntry{
		newEntry(filepath.Join(dir, "subdir")),
		newEntry(filepath.Join(dir, "file.txt")),
		newEntry(filepath.Join(dir, ".hidden")),
	}

	entry.SortChildren()

	if len(entry.children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(entry.children))
	}

	if entry.children[0].info.IsDir() {
		t.Error("expected first child after sort to be a file (lower mode value)")
	}

	if !entry.children[len(entry.children)-1].info.IsDir() {
		t.Error("expected last child after sort to be a directory (higher mode value)")
	}
}

// --- CleanChildList ---

func TestCleanChildList(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	origCount := len(entry.children)

	nonExistentChild := newEntry(filepath.Join(dir, "subdir", "nonexistent_cleanup_test"))
	entry.children = append(entry.children, nonExistentChild)

	if len(entry.children) != origCount+1 {
		t.Fatalf("expected %d children, got %d", origCount+1, len(entry.children))
	}

	entry.CleanChildListFunc((*SEntry).Exists)

	if len(entry.children) != origCount {
		t.Errorf("expected %d children after CleanChildListFunc, got %d", origCount, len(entry.children))
	}
}

func TestCleanChildListFunc_Custom(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(dir)
	entry.children = []*SEntry{
		newEntry(filepath.Join(dir, "file.txt")),
		newEntry(filepath.Join(dir, "subdir")),
		newEntry(filepath.Join(dir, ".hidden")),
	}

	entry.CleanChildListFunc(func(child *SEntry) bool {
		return !child.info.IsDir()
	})

	if len(entry.children) != 2 {
		t.Errorf("expected 2 children after removing dirs, got %d", len(entry.children))
	}
	for _, child := range entry.children {
		if child.info.IsDir() {
			t.Errorf("expected no directories after CleanChildListFunc")
		}
	}
}

// --- SetPathFilter / ApplyPathFilter ---

func TestSetPathFilter(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))

	filter := func(path string, info os.FileInfo) bool {
		return strings.HasSuffix(path, ".txt")
	}
	entry.SetPathFilter(filter)

	if entry.pathFilter == nil {
		t.Error("expected pathFilter to be set")
	}
	if entry.pathFilterCache == nil {
		t.Fatal("expected pathFilterCache to be initialized")
	}
	if len(entry.pathFilterCache) != 0 {
		t.Errorf("expected empty cache after SetPathFilter, got %d entries", len(entry.pathFilterCache))
	}
}

func TestApplyPathFilter(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(dir)
	entry.children = []*SEntry{
		newEntry(filepath.Join(dir, "file.txt")),
		newEntry(filepath.Join(dir, "subdir")),
	}

	entry.SetPathFilter(func(path string, info os.FileInfo) bool {
		return !info.IsDir()
	})

	entry.ApplyPathFilter()

	if len(entry.children) != 1 {
		t.Fatalf("expected 1 child after filter, got %d", len(entry.children))
	}
	if entry.children[0].info.Name() != "file.txt" {
		t.Errorf("expected remaining child to be 'file.txt', got %q", entry.children[0].info.Name())
	}
}

func TestApplyPathFilter_NilFilter(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(dir)
	entry.children = []*SEntry{
		newEntry(filepath.Join(dir, "file.txt")),
		newEntry(filepath.Join(dir, "subdir")),
	}

	entry.ApplyPathFilter()

	if len(entry.children) != 2 {
		t.Errorf("expected all 2 children to remain with nil filter, got %d", len(entry.children))
	}
}

func TestApplyPathFilter_CachesRejected(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(dir)
	entry.children = []*SEntry{
		newEntry(filepath.Join(dir, "file.txt")),
	}

	entry.SetPathFilter(func(path string, info os.FileInfo) bool {
		return false
	})

	entry.ApplyPathFilter()

	if len(entry.pathFilterCache) != 1 {
		t.Errorf("expected 1 path in cache after rejection, got %d", len(entry.pathFilterCache))
	}
	if len(entry.children) != 0 {
		t.Errorf("expected 0 children after all rejected, got %d", len(entry.children))
	}
}

// --- UpdateChildList ---

func TestUpdateChildList_NewEntries(t *testing.T) {
	dir := setupTestDir(t)
	subdir := filepath.Join(dir, "subdir")
	entry := newEntry(subdir)
	entry.UpdateChildList()

	mustWriteFile(t, filepath.Join(subdir, "newfile.txt"), []byte("new"), 0644)

	entry.UpdateChildList()

	if child := entry.GetChildByName("newfile.txt"); child == nil {
		t.Error("expected new child to be found after UpdateChildList")
	}
}

func TestUpdateChildList_OnFile(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	entry.UpdateChildList()

	if len(entry.children) != 0 {
		t.Errorf("expected 0 children for file entry, got %d", len(entry.children))
	}
}

func TestUpdateChildList_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "subdir")
	entry := newEntry(path)
	entry.UpdateChildList()

	if len(entry.children) != 2 {
		t.Fatalf("expected 2 children initially, got %d", len(entry.children))
	}

	t.Cleanup(func() { os.Chmod(path, 0755) })
	mustChmod(t, path, 0000)

	entry.UpdateChildList()

	if len(entry.children) != 2 {
		t.Errorf("expected 2 children preserved after permission error, got %d", len(entry.children))
	}
}

func TestUpdateChildList_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir", "empty"))
	entry.UpdateChildList()

	if len(entry.children) != 0 {
		t.Errorf("expected 0 children for empty dir, got %d", len(entry.children))
	}
}

func TestUpdateChildList_WithPathFilter(t *testing.T) {
	dir := setupTestDir(t)
	subdir := filepath.Join(dir, "subdir")
	entry := newEntry(subdir)

	entry.SetPathFilter(func(path string, info os.FileInfo) bool {
		return strings.HasSuffix(path, ".txt")
	})

	entry.UpdateChildList()

	if child := entry.GetChildByName("a.txt"); child == nil {
		t.Error("expected 'a.txt' to pass filter and be added as child")
	}
	if child := entry.GetChildByName("empty"); child != nil {
		t.Error("expected 'empty' to be rejected by filter")
	}
}

func TestUpdateChildList_WithPathFilter_CacheReuse(t *testing.T) {
	dir := setupTestDir(t)
	subdir := filepath.Join(dir, "subdir")
	entry := newEntry(subdir)

	var callCount int
	entry.SetPathFilter(func(path string, info os.FileInfo) bool {
		callCount++
		return false
	})

	entry.UpdateChildList()

	callCount = 0

	mustWriteFile(t, filepath.Join(subdir, "newfile.txt"), []byte("new"), 0644)

	entry.UpdateChildList()

	if callCount != 1 {
		t.Errorf("expected filter called 1 time (for newfile.txt only), got %d", callCount)
	}

	if !contains(entry.pathFilterCache, filepath.Join(subdir, "newfile.txt")) {
		t.Error("expected newfile.txt to be in cache after rejection")
	}
}

// --- GetChildByName ---

func TestGetChildByName_Found(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	child := entry.GetChildByName("a.txt")
	if child == nil {
		t.Fatal("expected to find child 'a.txt'")
	}
}

func TestGetChildByName_NotFound(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	if child := entry.GetChildByName("nonexistent.txt"); child != nil {
		t.Errorf("expected nil, got %v", child)
	}
}

func TestGetChildByName_EmptyName(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	if child := entry.GetChildByName(""); child != nil {
		t.Errorf("expected nil for empty name, got %v", child)
	}
}

// --- ForEachChild ---

func TestForEachChild(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	var count int
	entry.ForEachChild(func(child *SEntry) {
		count++
	})

	if count != 2 {
		t.Errorf("expected 2 children, got %d", count)
	}
}

func TestForEachChild_Empty(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir", "empty"))
	entry.UpdateChildList()

	var count int
	entry.ForEachChild(func(child *SEntry) {
		count++
	})

	if count != 0 {
		t.Errorf("expected 0 children for empty dir, got %d", count)
	}
}

// --- Recursive ---

func TestRecursive(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	var names []string
	entry.Recursive(func(child *SEntry) {
		names = append(names, child.info.Name())
	}, false)

	if !contains(names, "a.txt") {
		t.Errorf("expected 'a.txt' to be visited in recursive traversal, got %v", names)
	}
	if !contains(names, "empty") {
		t.Errorf("expected 'empty' to be visited in recursive traversal, got %v", names)
	}
}

func TestRecursive_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir", "empty"))
	entry.UpdateChildList()

	var count int
	entry.Recursive(func(child *SEntry) {
		count++
	}, false)

	if count != 0 {
		t.Errorf("expected 0 recursive calls for empty dir, got %d", count)
	}
}

func TestRecursive_IncludeSelf(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	var names []string
	entry.Recursive(func(n *SEntry) {
		names = append(names, n.info.Name())
	}, true)

	if !contains(names, "subdir") {
		t.Errorf("expected 'subdir' (self) to be visited with includeSelf=true, got %v", names)
	}
	if !contains(names, "a.txt") {
		t.Errorf("expected 'a.txt' to be visited, got %v", names)
	}
	if !contains(names, "empty") {
		t.Errorf("expected 'empty' to be visited, got %v", names)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 total visited entries, got %d: %v", len(names), names)
	}
}

// --- String ---

func TestString(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "file.txt"))

	result := entry.String()
	if !strings.Contains(result, "file.txt") {
		t.Errorf("expected String output to contain 'file.txt', got: %s", result)
	}
	parts := strings.Fields(result)
	if len(parts) < 3 {
		t.Errorf("expected at least 3 space-separated fields, got %d in: %s", len(parts), result)
	}
}

func TestString_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { entry.String() })
}

// --- StringRecursive ---

func TestStringRecursive(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir"))
	entry.UpdateChildList()

	result := entry.StringRecursive()
	if !strings.Contains(result, "a.txt") {
		t.Errorf("expected output to contain 'a.txt', got:\n%s", result)
	}
	if !strings.Contains(result, "empty") {
		t.Errorf("expected output to contain 'empty', got:\n%s", result)
	}
}

func TestStringRecursive_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "subdir", "empty"))
	entry.UpdateChildList()

	result := entry.StringRecursive()
	if result == "" {
		t.Error("expected non-empty output for empty dir, got empty string")
	}
}

func TestStringRecursive_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	entry := newEntry(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { entry.StringRecursive() })
}

// --- helpers ---

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
