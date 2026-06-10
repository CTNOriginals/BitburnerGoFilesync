package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	var baseDir = setupTestEnv()
	defer func() {
		os.Chmod(filepath.Join(baseDir, "noaccess_dir"), 0755)
		os.Chmod(filepath.Join(baseDir, "noaccess_file.txt"), 0644)
		os.RemoveAll(baseDir)
	}()

	testNewNode(baseDir)
	testGetPath(baseDir)
	testExists(baseDir)
	testIsModified(baseDir)
	testIsDirectory(baseDir)
	testChildren(baseDir)
	testSort(baseDir)
	testUpdate(baseDir)
	testTimeSinceModify(baseDir)
	testString(baseDir)
	testStringRecursive(baseDir)

	// return
	var node = newNode(config.Values.Directory)
	// var node = newNode(constants.WorkindDirectory)

	if node.infoError != nil {
		clog.Error(node.infoError)
	}

	node.UpdateChildList()
	// node.Recursive((*SNode).UpdateChildList)
	node.ForEachChild(func(child *SNode) {
		child.Recursive((*SNode).UpdateChildList)
	})

	clog.Infof("%s:\n%s", node.GetPath(), node.StringRecursive())
}

func setupTestEnv() string {
	baseDir, err := os.MkdirTemp("", "watcher-test-*")
	if err != nil {
		clog.Fatalf("Failed to create temp dir: %v", err)
	}

	os.WriteFile(filepath.Join(baseDir, "regular.txt"), []byte("hello"), 0644)

	os.MkdirAll(filepath.Join(baseDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(baseDir, "subdir", "nested.txt"), []byte("nested"), 0644)

	os.MkdirAll(filepath.Join(baseDir, "empty"), 0755)

	os.WriteFile(filepath.Join(baseDir, "noaccess_file.txt"), []byte("cant read me"), 0000)

	os.MkdirAll(filepath.Join(baseDir, "noaccess_dir"), 0755)
	os.WriteFile(filepath.Join(baseDir, "noaccess_dir", "trapped.txt"), []byte("trapped"), 0644)
	os.Chmod(filepath.Join(baseDir, "noaccess_dir"), 0000)

	os.WriteFile(filepath.Join(baseDir, "symlink_target.txt"), []byte("target"), 0644)
	os.Symlink("symlink_target.txt", filepath.Join(baseDir, "symlink_working"))
	os.Symlink("nonexistent_target", filepath.Join(baseDir, "symlink_broken"))
	os.Symlink("symlink_loop2", filepath.Join(baseDir, "symlink_loop1"))
	os.Symlink("symlink_loop1", filepath.Join(baseDir, "symlink_loop2"))

	os.WriteFile(filepath.Join(baseDir, "unicode_测试.txt"), []byte("unicode"), 0644)
	os.WriteFile(filepath.Join(baseDir, "spaces in name.txt"), []byte("spaces"), 0644)
	os.WriteFile(filepath.Join(baseDir, ".hidden"), []byte("hidden"), 0644)

	os.MkdirAll(filepath.Join(baseDir, "deep", "a", "b", "c", "d", "e"), 0755)
	os.WriteFile(filepath.Join(baseDir, "deep", "a", "b", "c", "d", "e", "leaf.txt"), []byte("deep"), 0644)

	return baseDir
}

func mkPath(baseDir string, elems ...string) string {
	return filepath.Join(append([]string{baseDir}, elems...)...)
}

func checkPASS(name string, ok bool) {
	if ok {
		clog.Infof("  PASS [%s]\n", name)
	} else {
		clog.Errorf("  FAIL [%s]\n", name)
	}
}

func checkPanic(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			clog.Infof("  BUG CONFIRMED [%s]: %v\n", name, r)
		} else {
			clog.Errorf("  BUG MISSED [%s]: expected panic but none occurred\n", name)
		}
	}()
	fn()
}

// -- newNode --

func testNewNode(baseDir string) {
	clog.Message("\n-- newNode --\n")

	checkPASS("existing file",
		func() bool {
			var n = newNode(mkPath(baseDir, "regular.txt"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("directory",
		func() bool {
			var n = newNode(mkPath(baseDir, "subdir"))
			return n.info != nil && n.infoError == nil && n.info.IsDir()
		}())

	checkPASS("empty directory",
		func() bool {
			var n = newNode(mkPath(baseDir, "empty"))
			return n.info != nil && n.infoError == nil && n.info.IsDir()
		}())

	checkPASS("working symlink",
		func() bool {
			var n = newNode(mkPath(baseDir, "symlink_working"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("unicode filename",
		func() bool {
			var n = newNode(mkPath(baseDir, "unicode_测试.txt"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("spaces in filename",
		func() bool {
			var n = newNode(mkPath(baseDir, "spaces in name.txt"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("hidden file",
		func() bool {
			var n = newNode(mkPath(baseDir, ".hidden"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("root path",
		func() bool {
			var n = newNode("/")
			return n.info != nil && n.infoError == nil && n.info.IsDir()
		}())

	checkPASS("nonexistent path",
		func() bool {
			var n = newNode(mkPath(baseDir, "nonexistent"))
			return n.info == nil && errors.Is(n.infoError, os.ErrNotExist)
		}())

	checkPASS("broken symlink",
		func() bool {
			var n = newNode(mkPath(baseDir, "symlink_broken"))
			return n.info == nil && errors.Is(n.infoError, os.ErrNotExist)
		}())

	checkPASS("symlink loop",
		func() bool {
			var n = newNode(mkPath(baseDir, "symlink_loop1"))
			return n.info == nil && errors.Is(n.infoError, syscall.ELOOP)
		}())

	checkPASS("trapped inside noaccess dir (EACCES)",
		func() bool {
			var n = newNode(mkPath(baseDir, "noaccess_dir", "trapped.txt"))
			return n.info == nil && errors.Is(n.infoError, os.ErrPermission)
		}())

	checkPASS("noaccess file (owner can stat)",
		func() bool {
			var n = newNode(mkPath(baseDir, "noaccess_file.txt"))
			return n.info != nil && n.infoError == nil
		}())

	checkPASS("path too long",
		func() bool {
			var n = newNode(mkPath(baseDir, strings.Repeat("a", 300)))
			return n.info == nil && errors.Is(n.infoError, syscall.ENAMETOOLONG)
		}())

	checkPASS("empty string path",
		func() bool {
			var n = newNode("")
			return n.info == nil && n.infoError != nil
		}())
}

// -- GetPath --

func testGetPath(baseDir string) {
	clog.Message("\n-- GetPath --\n")

	checkPASS("regular file",
		newNode(mkPath(baseDir, "regular.txt")).GetPath() == mkPath(baseDir, "regular.txt"))

	checkPASS("directory",
		newNode(mkPath(baseDir, "subdir")).GetPath() == mkPath(baseDir, "subdir"))

	checkPanic("nonexistent path", func() {
		newNode(mkPath(baseDir, "nonexistent")).GetPath()
	})

	checkPanic("broken symlink", func() {
		newNode(mkPath(baseDir, "symlink_broken")).GetPath()
	})

	checkPanic("symlink loop", func() {
		newNode(mkPath(baseDir, "symlink_loop1")).GetPath()
	})

	checkPanic("empty string path", func() {
		newNode("").GetPath()
	})
}

// -- Exists --

func testExists(baseDir string) {
	clog.Message("\n-- Exists --\n")

	checkPASS("valid file", newNode(mkPath(baseDir, "regular.txt")).Exists())
	checkPASS("directory", newNode(mkPath(baseDir, "subdir")).Exists())
	checkPASS("nonexistent", !newNode(mkPath(baseDir, "nonexistent")).Exists())
	checkPASS("broken symlink", !newNode(mkPath(baseDir, "symlink_broken")).Exists())
	checkPASS("symlink loop",
		newNode(mkPath(baseDir, "symlink_loop1")).Exists())
	checkPASS("trapped inside noaccess (EACCES)",
		newNode(mkPath(baseDir, "noaccess_dir", "trapped.txt")).Exists())
}

// -- IsModified --

func testIsModified(baseDir string) {
	clog.Message("\n-- IsModified --\n")

	var n = newNode(mkPath(baseDir, "regular.txt"))
	checkPASS("no change", !n.IsModified())

	os.WriteFile(mkPath(baseDir, "regular.txt"), []byte("modified"), 0644)
	checkPASS("content changed", n.IsModified())

	n.Update()
	checkPASS("after update no change", !n.IsModified())

	time.Sleep(time.Millisecond)
	os.Chtimes(mkPath(baseDir, "regular.txt"), time.Now(), time.Now())
	checkPASS("touch file (modtime changed)", n.IsModified())

	n.Update()
	os.Chmod(mkPath(baseDir, "regular.txt"), 0600)
	checkPASS("permission only (modtime unchanged)", !n.IsModified())

	var tmpPath = mkPath(baseDir, "_tmp_im_panic.txt")
	os.WriteFile(tmpPath, []byte("x"), 0644)
	var imNode = newNode(tmpPath)
	os.Remove(tmpPath)
	checkPanic("IsModified after deletion (getInfo returns nil)", func() {
		imNode.IsModified()
	})

	checkPanic("IsModified on nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).IsModified()
	})
}

// -- IsDirectory --

func testIsDirectory(baseDir string) {
	clog.Message("\n-- IsDirectory --\n")

	checkPASS("regular file", !newNode(mkPath(baseDir, "regular.txt")).IsDirectory())
	checkPASS("subdir", newNode(mkPath(baseDir, "subdir")).IsDirectory())
	checkPASS("empty dir", newNode(mkPath(baseDir, "empty")).IsDirectory())
	checkPanic("nonexistent path", func() {
		newNode(mkPath(baseDir, "nonexistent")).IsDirectory()
	})
}

// -- Children --

func testChildren(baseDir string) {
	clog.Message("\n-- Children --\n")

	var dirNode = newNode(mkPath(baseDir, "subdir"))
	dirNode.UpdateChildList()
	checkPASS("UpdateChildList populates children", len(dirNode.children) == 1)
	checkPASS("correct child name",
		dirNode.children[0].info.Name() == "nested.txt")

	checkPASS("GetChildByName found", dirNode.GetChildByName("nested.txt") != nil)
	checkPASS("GetChildByName not found", dirNode.GetChildByName("nosuch.txt") == nil)

	var count = 0
	dirNode.ForEachChild(func(child *SNode) {
		count++
	})
	checkPASS("ForEachChild count", count == 1)

	var recCount = 0
	dirNode.Recursive(func(child *SNode) {
		recCount++
	})
	checkPASS("Recursive count", recCount == 1)

	var emptyNode = newNode(mkPath(baseDir, "empty"))
	emptyNode.UpdateChildList()
	checkPASS("empty directory", len(emptyNode.children) == 0)

	var fileNode = newNode(mkPath(baseDir, "regular.txt"))
	fileNode.UpdateChildList()
	checkPASS("file UpdateChildList is no-op", len(fileNode.children) == 0)

	var noaccNode = newNode(mkPath(baseDir, "noaccess_dir"))
	noaccNode.UpdateChildList()
	checkPASS("noaccess dir (EACCES on ReadDir)", len(noaccNode.children) == 0)

	os.WriteFile(mkPath(baseDir, "subdir", "newfile.txt"), []byte("new"), 0644)
	dirNode.UpdateChildList()
	checkPASS("added child appears", len(dirNode.children) == 2)

	os.Remove(mkPath(baseDir, "subdir", "newfile.txt"))
	dirNode.CleanChildList()
	checkPASS("CleanChildList removes deleted child", len(dirNode.children) == 1)

	var deepNode = newNode(mkPath(baseDir, "deep"))
	deepNode.UpdateChildList()
	var deepCount = 0
	deepNode.Recursive(func(child *SNode) {
		deepCount++
	})
	checkPASS("deep nesting Recursive count", deepCount == 6)

	checkPanic("UpdateChildList on nonexistent node (GetPath panics)", func() {
		newNode(mkPath(baseDir, "nonexistent")).UpdateChildList()
	})

	var parentWithNilChild = newNode(mkPath(baseDir, "empty"))
	var nilInfoChild = newNode(mkPath(baseDir, "symlink_broken"))
	parentWithNilChild.children = append(parentWithNilChild.children, nilInfoChild)
	checkPanic("GetChildByName on child with nil info", func() {
		parentWithNilChild.GetChildByName("anything")
	})
}

// -- Sort --

func testSort(baseDir string) {
	clog.Message("\n-- Sort --\n")

	var sortDirPath = mkPath(baseDir, "sort_test")
	os.MkdirAll(sortDirPath, 0755)
	defer os.RemoveAll(sortDirPath)

	os.WriteFile(mkPath(sortDirPath, "z_file.txt"), []byte("z"), 0644)
	os.WriteFile(mkPath(sortDirPath, "a_file.txt"), []byte("a"), 0644)
	os.MkdirAll(mkPath(sortDirPath, "m_dir"), 0755)

	var sortNode = newNode(sortDirPath)
	sortNode.UpdateChildList()
	sortNode.SortChildren()

	var names []string
	for _, c := range sortNode.children {
		names = append(names, c.info.Name())
	}

	checkPASS("files before dirs, alphabetical",
		len(names) == 3 &&
			names[0] == "a_file.txt" &&
			names[1] == "z_file.txt" &&
			names[2] == "m_dir")

	var badSort = newNode(sortDirPath)
	var goodChild = newNode(mkPath(baseDir, "regular.txt"))
	var nilChild = newNode(mkPath(baseDir, "nonexistent"))
	badSort.children = append(badSort.children, goodChild, nilChild)
	checkPanic("SortChildren with nil-info child", func() {
		badSort.SortChildren()
	})
}

// -- Update --

func testUpdate(baseDir string) {
	clog.Message("\n-- Update --\n")

	var n = newNode(mkPath(baseDir, "regular.txt"))
	n.Update()
	checkPASS("update valid file", n.info != nil && n.infoError == nil)

	var tmpPath = mkPath(baseDir, "_tmp_update.txt")
	os.WriteFile(tmpPath, []byte("x"), 0644)
	var dn = newNode(tmpPath)
	os.Remove(tmpPath)
	dn.Update()
	checkPASS("update after deletion",
		dn.info == nil && errors.Is(dn.infoError, os.ErrNotExist))

	checkPanic("Update on nil-info node (GetPath panics)", func() {
		newNode(mkPath(baseDir, "nonexistent")).Update()
	})
}

// -- GetTimeSinceModify --

func testTimeSinceModify(baseDir string) {
	clog.Message("\n-- GetTimeSinceModify --\n")

	checkPASS("valid file", newNode(mkPath(baseDir, "regular.txt")).GetTimeSinceModify() > 0)

	checkPanic("nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).GetTimeSinceModify()
	})
}

// -- String --

func testString(baseDir string) {
	clog.Message("\n-- String --\n")

	checkPASS("valid file", strings.Contains(newNode(mkPath(baseDir, "regular.txt")).String(), "regular.txt"))
	checkPASS("directory", strings.Contains(newNode(mkPath(baseDir, "subdir")).String(), "subdir"))
	checkPanic("nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).String()
	})
}

// -- StringRecursive --

func testStringRecursive(baseDir string) {
	clog.Message("\n-- StringRecursive --\n")

	var dirNode = newNode(mkPath(baseDir, "subdir"))
	dirNode.UpdateChildList()
	clog.Infof("  subdir tree:\n%s\n", dirNode.StringRecursive())
	checkPASS("subdir tree", strings.Contains(dirNode.StringRecursive(), "nested.txt"))

	var deepNode = newNode(mkPath(baseDir, "deep"))
	deepNode.UpdateChildList()
	clog.Infof("  deep tree:\n%s\n", deepNode.StringRecursive())
	checkPASS("deep tree", strings.Contains(deepNode.StringRecursive(), "leaf.txt"))
}
