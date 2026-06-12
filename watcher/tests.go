package watcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	var baseDir = setupTestEnv()
	defer func() {
		os.Chmod(filepath.Join(baseDir, "noaccess_dir"), 0755)
		os.Chmod(filepath.Join(baseDir, "noaccess_file.txt"), 0644)
		os.RemoveAll(baseDir)
	}()

	var grand sTestRun
	grand.init("grand")

	runSection("newNode", testNewNode, &grand, baseDir)
	runSection("GetPath", testGetPath, &grand, baseDir)
	runSection("Exists", testExists, &grand, baseDir)
	runSection("IsModified", testIsModified, &grand, baseDir)
	runSection("IsDirectory", testIsDirectory, &grand, baseDir)
	runSection("Children", testChildren, &grand, baseDir)
	runSection("Sort", testSort, &grand, baseDir)
	runSection("Update", testUpdate, &grand, baseDir)
	runSection("TimeSinceModify", testTimeSinceModify, &grand, baseDir)
	runSection("String", testString, &grand, baseDir)
	runSection("StringRecursive", testStringRecursive, &grand, baseDir)

	clog.Messagef("\n========== Grand Total: %d passed, %d failed, %d panicked ==========",
		grand.passed, grand.failed, grand.panicked)
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

type sTestRun struct {
	section  string
	passed   int
	failed   int
	panicked int
}

func (tr *sTestRun) init(section string) {
	tr.section = section
	tr.passed = 0
	tr.failed = 0
	tr.panicked = 0
}

func (tr *sTestRun) check(desc string, ok bool, details ...string) {
	if ok {
		clog.Messagef("  PASS  %s", desc)
		tr.passed++
	} else {
		var detail = strings.Join(details, ", ")
		if detail != "" {
			clog.Messagef("  FAIL  %s  << %s >>", desc, detail)
		} else {
			clog.Messagef("  FAIL  %s", desc)
		}
		tr.failed++
	}
}

func (tr *sTestRun) expectPanic(desc string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			clog.Messagef("  PASS  %s  (panic: %v)", desc, r)
			tr.passed++
		} else {
			clog.Messagef("  FAIL  %s  (expected panic, none occurred)", desc)
			tr.failed++
		}
	}()
	fn()
}

func (tr *sTestRun) summary() {
	clog.Messagef("-- %s: %d passed, %d failed, %d panicked --",
		tr.section, tr.passed, tr.failed, tr.panicked)
}

func runSection(section string, fn func(*sTestRun, string), grand *sTestRun, baseDir string) {
	var tr sTestRun
	tr.init(section)
	clog.Messagef("\n-- %s --\n", section)
	fn(&tr, baseDir)
	tr.summary()
	grand.passed += tr.passed
	grand.failed += tr.failed
	grand.panicked += tr.panicked
}

// -- newNode --

func testNewNode(tr *sTestRun, baseDir string) {
	var regN = newNode(mkPath(baseDir, "regular.txt"))
	tr.check("existing file",
		regN.info != nil && regN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", regN.info, regN.infoError))

	var subdirN = newNode(mkPath(baseDir, "subdir"))
	tr.check("directory",
		subdirN.info != nil && subdirN.infoError == nil && subdirN.info.IsDir(),
		fmt.Sprintf("info=%v, infoError=%v", subdirN.info, subdirN.infoError))

	var emptyN = newNode(mkPath(baseDir, "empty"))
	tr.check("empty directory",
		emptyN.info != nil && emptyN.infoError == nil && emptyN.info.IsDir(),
		fmt.Sprintf("info=%v, infoError=%v", emptyN.info, emptyN.infoError))

	var symwkN = newNode(mkPath(baseDir, "symlink_working"))
	tr.check("working symlink",
		symwkN.info != nil && symwkN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", symwkN.info, symwkN.infoError))

	var utfN = newNode(mkPath(baseDir, "unicode_测试.txt"))
	tr.check("unicode filename",
		utfN.info != nil && utfN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", utfN.info, utfN.infoError))

	var spacesN = newNode(mkPath(baseDir, "spaces in name.txt"))
	tr.check("spaces in filename",
		spacesN.info != nil && spacesN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", spacesN.info, spacesN.infoError))

	var hiddenN = newNode(mkPath(baseDir, ".hidden"))
	tr.check("hidden file",
		hiddenN.info != nil && hiddenN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", hiddenN.info, hiddenN.infoError))

	var rootN = newNode("/")
	tr.check("root path",
		rootN.info != nil && rootN.infoError == nil && rootN.info.IsDir(),
		fmt.Sprintf("info=%v, infoError=%v", rootN.info, rootN.infoError))

	var nonexistN = newNode(mkPath(baseDir, "nonexistent"))
	tr.check("nonexistent path",
		nonexistN.info == nil && errors.Is(nonexistN.infoError, os.ErrNotExist),
		fmt.Sprintf("info=%v, infoError=%v", nonexistN.info, nonexistN.infoError))

	var brokenN = newNode(mkPath(baseDir, "symlink_broken"))
	tr.check("broken symlink",
		brokenN.info == nil && errors.Is(brokenN.infoError, os.ErrNotExist),
		fmt.Sprintf("info=%v, infoError=%v", brokenN.info, brokenN.infoError))

	var loopN = newNode(mkPath(baseDir, "symlink_loop1"))
	tr.check("symlink loop",
		loopN.info == nil && errors.Is(loopN.infoError, syscall.ELOOP),
		fmt.Sprintf("info=%v, infoError=%v", loopN.info, loopN.infoError))

	var trappedN = newNode(mkPath(baseDir, "noaccess_dir", "trapped.txt"))
	tr.check("trapped inside noaccess dir (EACCES)",
		trappedN.info == nil && errors.Is(trappedN.infoError, os.ErrPermission),
		fmt.Sprintf("info=%v, infoError=%v", trappedN.info, trappedN.infoError))

	var noaccfileN = newNode(mkPath(baseDir, "noaccess_file.txt"))
	tr.check("noaccess file (owner can stat)",
		noaccfileN.info != nil && noaccfileN.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", noaccfileN.info, noaccfileN.infoError))

	var longN = newNode(mkPath(baseDir, strings.Repeat("a", 300)))
	tr.check("path too long",
		longN.info == nil && errors.Is(longN.infoError, syscall.ENAMETOOLONG),
		fmt.Sprintf("info=%v, infoError=%v", longN.info, longN.infoError))

	var emptysN = newNode("")
	tr.check("empty string path",
		emptysN.info == nil && emptysN.infoError != nil,
		fmt.Sprintf("info=%v, infoError=%v", emptysN.info, emptysN.infoError))
}

// -- GetPath --

func testGetPath(tr *sTestRun, baseDir string) {
	tr.check("regular file",
		newNode(mkPath(baseDir, "regular.txt")).GetPath() == mkPath(baseDir, "regular.txt"))

	tr.check("directory",
		newNode(mkPath(baseDir, "subdir")).GetPath() == mkPath(baseDir, "subdir"))

	tr.expectPanic("nonexistent path", func() {
		newNode(mkPath(baseDir, "nonexistent")).GetPath()
	})

	tr.expectPanic("broken symlink", func() {
		newNode(mkPath(baseDir, "symlink_broken")).GetPath()
	})

	tr.expectPanic("symlink loop", func() {
		newNode(mkPath(baseDir, "symlink_loop1")).GetPath()
	})

	tr.expectPanic("empty string path", func() {
		newNode("").GetPath()
	})
}

// -- Exists --

func testExists(tr *sTestRun, baseDir string) {
	tr.check("valid file", newNode(mkPath(baseDir, "regular.txt")).Exists())
	tr.check("directory", newNode(mkPath(baseDir, "subdir")).Exists())
	tr.check("nonexistent", !newNode(mkPath(baseDir, "nonexistent")).Exists())
	tr.check("broken symlink", !newNode(mkPath(baseDir, "symlink_broken")).Exists())
	tr.check("symlink loop",
		newNode(mkPath(baseDir, "symlink_loop1")).Exists())
	tr.check("trapped inside noaccess (EACCES)",
		newNode(mkPath(baseDir, "noaccess_dir", "trapped.txt")).Exists())
}

// -- IsModified --

func testIsModified(tr *sTestRun, baseDir string) {
	var n = newNode(mkPath(baseDir, "regular.txt"))
	tr.check("no change", !n.IsModified())

	os.WriteFile(mkPath(baseDir, "regular.txt"), []byte("modified"), 0644)
	tr.check("content changed", n.IsModified())

	n.Update()
	tr.check("after update no change", !n.IsModified())

	time.Sleep(time.Millisecond)
	os.Chtimes(mkPath(baseDir, "regular.txt"), time.Now(), time.Now())
	tr.check("touch file (modtime changed)", n.IsModified())

	n.Update()
	os.Chmod(mkPath(baseDir, "regular.txt"), 0600)
	tr.check("permission only (modtime unchanged)", !n.IsModified())

	var tmpPath = mkPath(baseDir, "_tmp_im_panic.txt")
	os.WriteFile(tmpPath, []byte("x"), 0644)
	var imNode = newNode(tmpPath)
	os.Remove(tmpPath)
	tr.expectPanic("IsModified after deletion (getInfo returns nil)", func() {
		imNode.IsModified()
	})

	tr.expectPanic("IsModified on nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).IsModified()
	})
}

// -- IsDirectory --

func testIsDirectory(tr *sTestRun, baseDir string) {
	tr.check("regular file", !newNode(mkPath(baseDir, "regular.txt")).IsDirectory())
	tr.check("subdir", newNode(mkPath(baseDir, "subdir")).IsDirectory())
	tr.check("empty dir", newNode(mkPath(baseDir, "empty")).IsDirectory())
	tr.expectPanic("nonexistent path", func() {
		newNode(mkPath(baseDir, "nonexistent")).IsDirectory()
	})
}

// -- Children --

func testChildren(tr *sTestRun, baseDir string) {
	var dirNode = newNode(mkPath(baseDir, "subdir"))
	dirNode.UpdateChildList()
	tr.check("UpdateChildList populates children", len(dirNode.children) == 1)
	tr.check("correct child name",
		dirNode.children[0].info.Name() == "nested.txt")

	tr.check("GetChildByName found", dirNode.GetChildByName("nested.txt") != nil)
	tr.check("GetChildByName not found", dirNode.GetChildByName("nosuch.txt") == nil)

	var count = 0
	dirNode.ForEachChild(func(child *SNode) {
		count++
	})
	tr.check("ForEachChild count", count == 1)

	var recCount = 0
	dirNode.Recursive(func(child *SNode) {
		recCount++
	})
	tr.check("Recursive count", recCount == 1)

	var emptyNode = newNode(mkPath(baseDir, "empty"))
	emptyNode.UpdateChildList()
	tr.check("empty directory", len(emptyNode.children) == 0)

	var fileNode = newNode(mkPath(baseDir, "regular.txt"))
	fileNode.UpdateChildList()
	tr.check("file UpdateChildList is no-op", len(fileNode.children) == 0)

	var noaccNode = newNode(mkPath(baseDir, "noaccess_dir"))
	noaccNode.UpdateChildList()
	tr.check("noaccess dir (EACCES on ReadDir)", len(noaccNode.children) == 0)

	os.WriteFile(mkPath(baseDir, "subdir", "newfile.txt"), []byte("new"), 0644)
	dirNode.UpdateChildList()
	tr.check("added child appears", len(dirNode.children) == 2)

	os.Remove(mkPath(baseDir, "subdir", "newfile.txt"))
	// dirNode.Recursive((*SNode).Update)
	dirNode.CleanChildList()
	tr.check("CleanChildList removes deleted child", len(dirNode.children) == 1,
		fmt.Sprintf("len=%d", len(dirNode.children)))

	var deepNode = newNode(mkPath(baseDir, "deep"))
	deepNode.UpdateChildList()
	var deepCount = 0
	deepNode.Recursive(func(child *SNode) {
		deepCount++
	})
	tr.check("deep nesting Recursive count", deepCount == 6)

	tr.expectPanic("UpdateChildList on nonexistent node (GetPath panics)", func() {
		newNode(mkPath(baseDir, "nonexistent")).UpdateChildList()
	})

	var parentWithNilChild = newNode(mkPath(baseDir, "empty"))
	var nilInfoChild = newNode(mkPath(baseDir, "symlink_broken"))
	parentWithNilChild.children = append(parentWithNilChild.children, nilInfoChild)
	tr.expectPanic("GetChildByName on child with nil info", func() {
		parentWithNilChild.GetChildByName("anything")
	})
}

// -- Sort --

func testSort(tr *sTestRun, baseDir string) {
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

	tr.check("files before dirs, alphabetical",
		len(names) == 3 &&
			names[0] == "a_file.txt" &&
			names[1] == "z_file.txt" &&
			names[2] == "m_dir")

	var badSort = newNode(sortDirPath)
	var goodChild = newNode(mkPath(baseDir, "regular.txt"))
	var nilChild = newNode(mkPath(baseDir, "nonexistent"))
	badSort.children = append(badSort.children, goodChild, nilChild)
	tr.expectPanic("SortChildren with nil-info child", func() {
		badSort.SortChildren()
	})
}

// -- Update --

func testUpdate(tr *sTestRun, baseDir string) {
	var n = newNode(mkPath(baseDir, "regular.txt"))
	n.Update()
	tr.check("update valid file", n.info != nil && n.infoError == nil,
		fmt.Sprintf("info=%v, infoError=%v", n.info, n.infoError))

	var tmpPath = mkPath(baseDir, "_tmp_update.txt")
	os.WriteFile(tmpPath, []byte("x"), 0644)
	var dn = newNode(tmpPath)
	os.Remove(tmpPath)
	dn.Update()
	tr.check("update after deletion",
		dn.info == nil && errors.Is(dn.infoError, os.ErrNotExist),
		fmt.Sprintf("info=%v, infoError=%v", dn.info, dn.infoError))

	tr.expectPanic("Update on nil-info node (GetPath panics)", func() {
		newNode(mkPath(baseDir, "nonexistent")).Update()
	})
}

// -- GetTimeSinceModify --

func testTimeSinceModify(tr *sTestRun, baseDir string) {
	tr.check("valid file", newNode(mkPath(baseDir, "regular.txt")).GetTimeSinceModify() > 0)

	tr.expectPanic("nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).GetTimeSinceModify()
	})
}

// -- String --

func testString(tr *sTestRun, baseDir string) {
	tr.check("valid file", strings.Contains(newNode(mkPath(baseDir, "regular.txt")).String(), "regular.txt"))
	tr.check("directory", strings.Contains(newNode(mkPath(baseDir, "subdir")).String(), "subdir"))
	tr.expectPanic("nonexistent node", func() {
		newNode(mkPath(baseDir, "nonexistent")).String()
	})
}

// -- StringRecursive --

func testStringRecursive(tr *sTestRun, baseDir string) {
	var dirNode = newNode(mkPath(baseDir, "subdir"))
	dirNode.UpdateChildList()
	clog.Infof("  subdir tree:\n%s\n", dirNode.StringRecursive())
	tr.check("subdir tree", strings.Contains(dirNode.StringRecursive(), "nested.txt"))

	var deepNode = newNode(mkPath(baseDir, "deep"))
	deepNode.UpdateChildList()
	clog.Infof("  deep tree:\n%s\n", deepNode.StringRecursive())
	tr.check("deep tree", strings.Contains(deepNode.StringRecursive(), "leaf.txt"))
}
