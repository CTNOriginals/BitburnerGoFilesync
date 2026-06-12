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

// --- newNode ---

func TestNewNode_File(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	if node.info == nil {
		t.Fatal("expected non-nil info for existing file")
	}
	if node.infoError != nil {
		t.Fatalf("expected nil error, got %v", node.infoError)
	}
	if node.info.Name() != "file.txt" {
		t.Errorf("expected name 'file.txt', got %q", node.info.Name())
	}
}

func TestNewNode_Directory(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))

	if node.info == nil {
		t.Fatal("expected non-nil info for existing directory")
	}
	if node.infoError != nil {
		t.Fatalf("expected nil error, got %v", node.infoError)
	}
	if !node.info.IsDir() {
		t.Error("expected IsDir to be true")
	}
}

func TestNewNode_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	if node.info != nil {
		t.Error("expected nil info for non-existent path")
	}
	if node.infoError == nil {
		t.Fatal("expected error for non-existent path")
	}
	if _, ok := node.infoError.(*os.PathError); !ok {
		t.Errorf("expected *os.PathError, got %T", node.infoError)
	}
	if !os.IsNotExist(node.infoError) {
		t.Errorf("expected not-exist error, got %v", node.infoError)
	}
}

func TestNewNode_NotADirectory(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "notdir_test", "child")
	node := newNode(path)

	if node.info != nil {
		t.Error("expected nil info for ENOTDIR path")
	}
	if node.infoError == nil {
		t.Fatal("expected error for ENOTDIR path")
	}
}

func TestNewNode_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	noperm := filepath.Join(dir, "noperm")
	t.Cleanup(func() { os.Chmod(noperm, 0755) })
	mustChmod(t, noperm, 0000)

	path := filepath.Join(noperm, "secret.txt")
	node := newNode(path)

	if node.info != nil {
		t.Error("expected nil info for permission-denied path")
	}
	if node.infoError == nil {
		t.Fatal("expected error for permission-denied path")
	}
	if pathErr, ok := node.infoError.(*os.PathError); ok {
		if !os.IsPermission(pathErr) {
			t.Errorf("expected permission error, got %v", pathErr.Err)
		}
	} else {
		t.Errorf("expected *os.PathError, got %T", node.infoError)
	}
}

func TestNewNode_SymlinkLoop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink loop test requires unix")
	}
	dir := setupTestDir(t)
	path := filepath.Join(dir, "chain", "link1")
	node := newNode(path)

	if node.info != nil {
		t.Error("expected nil info for symlink loop path")
	}
	if node.infoError == nil {
		t.Fatal("expected error for symlink loop path")
	}
}

func TestNewNode_PathTooLong(t *testing.T) {
	dir := setupTestDir(t)
	longPath := dir
	for len(longPath) < 4096 {
		longPath = filepath.Join(longPath, "a")
	}
	node := newNode(longPath)

	if node.info != nil {
		t.Error("expected nil info for path-too-long")
	}
	if node.infoError == nil {
		t.Fatal("expected error for path-too-long")
	}
}

// --- Exists ---

func TestExists_Existing(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	if !node.Exists() {
		t.Error("expected Exists() to be true for existing file")
	}
}

func TestExists_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	if node.Exists() {
		t.Error("expected Exists() to be false for non-existent path")
	}
}

// --- IsDirectory ---

func TestIsDirectory_File(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	if node.IsDirectory() {
		t.Error("expected IsDirectory to be false for a file")
	}
}

func TestIsDirectory_Directory(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))

	if !node.IsDirectory() {
		t.Error("expected IsDirectory to be true for a directory")
	}
}

func TestIsDirectory_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { node.IsDirectory() })
}

// --- GetPath ---

func TestGetPath(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	expected := filepath.Join(dir, "file.txt")
	if got := node.GetPath(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGetPath_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { node.GetPath() })
}

// --- IsModified ---

func TestIsModified_Unchanged(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	if node.IsModified() {
		t.Error("expected IsModified to be false for unchanged file")
	}
}

func TestIsModified_Changed(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	node := newNode(path)

	time.Sleep(10 * time.Millisecond)

	mustWriteFile(t, path, []byte("modified"), 0644)

	if !node.IsModified() {
		t.Error("expected IsModified to be true after file modification")
	}
}

func TestIsModified_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { node.IsModified() })
}

// --- Update ---

func TestUpdate_NoChange(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	node.Update()

	if node.infoError != nil {
		t.Errorf("expected nil infoError after Update on unchanged file, got %v", node.infoError)
	}
}

func TestUpdate_AfterModify(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	node := newNode(path)

	origModTime := node.info.ModTime()

	time.Sleep(10 * time.Millisecond)
	mustWriteFile(t, path, []byte("modified"), 0644)

	node.Update()

	if node.info.ModTime() == origModTime {
		t.Error("expected ModTime to change after Update")
	}
	if node.infoError != nil {
		t.Errorf("expected nil infoError, got %v", node.infoError)
	}
}

func TestUpdate_AfterDelete(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "file.txt")
	node := newNode(path)

	origName := node.info.Name()

	os.Remove(path)

	node.Update()

	if node.infoError == nil {
		t.Error("expected infoError after file deletion")
	}
	if !os.IsNotExist(node.infoError) {
		t.Errorf("expected not-exist error, got %v", node.infoError)
	}
	if node.info == nil {
		t.Fatal("expected info to be preserved after Update with error")
	}
	if node.info.Name() != origName {
		t.Errorf("expected name %q to be preserved, got %q", origName, node.info.Name())
	}
}

func TestUpdate_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	noperm := filepath.Join(dir, "noperm")
	path := filepath.Join(noperm, "secret.txt")
	node := newNode(path)

	origName := node.info.Name()

	t.Cleanup(func() { os.Chmod(noperm, 0755) })
	mustChmod(t, noperm, 0000)

	node.Update()

	if node.infoError == nil {
		t.Error("expected infoError after permission denied")
	}
	if node.info == nil {
		t.Fatal("expected info to be preserved")
	}
	if node.info.Name() != origName {
		t.Errorf("expected name %q to be preserved, got %q", origName, node.info.Name())
	}
}

// --- SortChildren ---

func TestSortChildren(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(dir)
	node.children = []*SNode{
		newNode(filepath.Join(dir, "subdir")),
		newNode(filepath.Join(dir, "file.txt")),
		newNode(filepath.Join(dir, ".hidden")),
	}

	node.SortChildren()

	if len(node.children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(node.children))
	}

	if node.children[0].info.IsDir() {
		t.Error("expected first child after sort to be a file (lower mode value)")
	}

	if !node.children[len(node.children)-1].info.IsDir() {
		t.Error("expected last child after sort to be a directory (higher mode value)")
	}
}

// --- CleanChildList ---

func TestCleanChildList(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	origCount := len(node.children)

	nonExistentChild := newNode(filepath.Join(dir, "subdir", "nonexistent_cleanup_test"))
	node.children = append(node.children, nonExistentChild)

	if len(node.children) != origCount+1 {
		t.Fatalf("expected %d children, got %d", origCount+1, len(node.children))
	}

	node.CleanChildList()

	if len(node.children) != origCount {
		t.Errorf("expected %d children after CleanChildList, got %d", origCount, len(node.children))
	}
}

// --- UpdateChildList ---

func TestUpdateChildList_NewEntries(t *testing.T) {
	dir := setupTestDir(t)
	subdir := filepath.Join(dir, "subdir")
	node := newNode(subdir)
	node.UpdateChildList()

	mustWriteFile(t, filepath.Join(subdir, "newfile.txt"), []byte("new"), 0644)

	node.UpdateChildList()

	if child := node.GetChildByName("newfile.txt"); child == nil {
		t.Error("expected new child to be found after UpdateChildList")
	}
}

func TestUpdateChildList_OnFile(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	node.UpdateChildList()

	if len(node.children) != 0 {
		t.Errorf("expected 0 children for file node, got %d", len(node.children))
	}
}

func TestUpdateChildList_PermissionDenied(t *testing.T) {
	dir := setupTestDir(t)
	path := filepath.Join(dir, "subdir")
	node := newNode(path)
	node.UpdateChildList()

	if len(node.children) != 2 {
		t.Fatalf("expected 2 children initially, got %d", len(node.children))
	}

	t.Cleanup(func() { os.Chmod(path, 0755) })
	mustChmod(t, path, 0000)

	node.UpdateChildList()

	if len(node.children) != 2 {
		t.Errorf("expected 2 children preserved after permission error, got %d", len(node.children))
	}
}

func TestUpdateChildList_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir", "empty"))
	node.UpdateChildList()

	if len(node.children) != 0 {
		t.Errorf("expected 0 children for empty dir, got %d", len(node.children))
	}
}

// --- GetChildByName ---

func TestGetChildByName_Found(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	child := node.GetChildByName("a.txt")
	if child == nil {
		t.Fatal("expected to find child 'a.txt'")
	}
}

func TestGetChildByName_NotFound(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	if child := node.GetChildByName("nonexistent.txt"); child != nil {
		t.Errorf("expected nil, got %v", child)
	}
}

func TestGetChildByName_EmptyName(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	if child := node.GetChildByName(""); child != nil {
		t.Errorf("expected nil for empty name, got %v", child)
	}
}

// --- ForEachChild ---

func TestForEachChild(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	var count int
	node.ForEachChild(func(child *SNode) {
		count++
	})

	if count != 2 {
		t.Errorf("expected 2 children, got %d", count)
	}
}

func TestForEachChild_Empty(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir", "empty"))
	node.UpdateChildList()

	var count int
	node.ForEachChild(func(child *SNode) {
		count++
	})

	if count != 0 {
		t.Errorf("expected 0 children for empty dir, got %d", count)
	}
}

// --- Recursive ---

func TestRecursive(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	var names []string
	node.Recursive(func(child *SNode) {
		names = append(names, child.info.Name())
	})

	if !contains(names, "a.txt") {
		t.Errorf("expected 'a.txt' to be visited in recursive traversal, got %v", names)
	}
	if !contains(names, "empty") {
		t.Errorf("expected 'empty' to be visited in recursive traversal, got %v", names)
	}
}

func TestRecursive_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir", "empty"))
	node.UpdateChildList()

	var count int
	node.Recursive(func(child *SNode) {
		count++
	})

	if count != 0 {
		t.Errorf("expected 0 recursive calls for empty dir, got %d", count)
	}
}

// --- String ---

func TestString(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "file.txt"))

	result := node.String()
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
	node := newNode(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { node.String() })
}

// --- StringRecursive ---

func TestStringRecursive(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir"))
	node.UpdateChildList()

	result := node.StringRecursive()
	if !strings.Contains(result, "a.txt") {
		t.Errorf("expected output to contain 'a.txt', got:\n%s", result)
	}
	if !strings.Contains(result, "empty") {
		t.Errorf("expected output to contain 'empty', got:\n%s", result)
	}
}

func TestStringRecursive_EmptyDir(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "subdir", "empty"))
	node.UpdateChildList()

	result := node.StringRecursive()
	if result == "" {
		t.Error("expected non-empty output for empty dir, got empty string")
	}
}

func TestStringRecursive_NonExistent(t *testing.T) {
	dir := setupTestDir(t)
	node := newNode(filepath.Join(dir, "does_not_exist"))

	assertPanic(t, func() { node.StringRecursive() })
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
