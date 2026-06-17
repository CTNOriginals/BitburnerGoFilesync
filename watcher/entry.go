package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

type FPathFilter func(path string, info os.FileInfo) bool

/*
SEntry can hold info about a file or directory.
It essentially is a wrapper for the os package.

It does not keep track of itself or its children in any way,
the stored information remains unchanged until [SEntry.Update] is called.

When creating a [newEntry], the entry will not populate its children,
it is instead expected that [SEntry.UpdateChildList] is called after creation.
The only exception to this is when [SEntry.UpdateChildList] creates a [newEntry] internally,
it will also call [SEntry.UpdateChildList] on this new entry.

Most functions do not recurse into children,
for this to happen, you may use [SEntry.ForEachChild] or [SEntry.Recursive],
and pass in the function that should be ran on children.
Example:

	SEntry.Recursive((*SEntry).Update)
*/
type SEntry struct {
	// The directory path that contains this entry.
	dir string

	// The last know info of this entry.
	info os.FileInfo
	// If an error occured while trying to get the info,
	// it will be stored in here.
	infoError error

	children []*SEntry

	// A function that runs when a new child is discovered.
	// The child will only be added to children if pathFilterFn returns true.
	// All children of this entry will inherit this function recursively.
	pathFilter      FPathFilter
	pathFilterCache []string
}

// newEntry Creates a new entry from path.
// It does not pre-populate the entry's list of children.
//
// If an error occurs while getting the entries info,
// the error will be stored in this entry so that it
// can be evaluated with the provided functions.
func newEntry(path string) *SEntry {
	var info, err = os.Stat(path)

	var entry = SEntry{
		dir:       filepath.Dir(path),
		info:      info,
		infoError: err,
	}

	return &entry
}

func (this SEntry) GetPath() string {
	return filepath.Join(this.dir, this.info.Name())
}

func (this SEntry) getInfo() (os.FileInfo, error) {
	var info, err = os.Stat(this.GetPath())

	if err != nil {
		var pathError, isPathError = err.(*os.PathError)

		if !isPathError {
			clog.Fatalf("Received unknown error:\n%v", err)
			runtime.Goexit()
		}

		return info, pathError
	}

	return info, nil
}

func (this SEntry) Exists() bool {
	return !os.IsNotExist(this.infoError)
}

func (this SEntry) IsModified() bool {
	var info, _ = this.getInfo()
	return this.info.ModTime() != info.ModTime()
}

func (this SEntry) IsDirectory() bool {
	return this.info.IsDir()
}

func (this SEntry) GetChildDirectories() []*SEntry {
	var dirs = make([]*SEntry, 0)

	for _, child := range this.children {
		if child.IsDirectory() {
			dirs = append(dirs, child)
		}
	}

	return dirs
}

func (this SEntry) GetChildByName(name string) *SEntry {
	for _, child := range this.children {
		if child.info.Name() == name {
			return child
		}
	}

	return nil
}

// Sets the pathFilter of this entry and clear pathFilterCache.
// Does NOT apply this filter to existing children,
// to do that, run this function in [SEntry.Recursive].
//
// To re-apply this pathFilter to the current list of children,
// run [SEntry.ApplyPathFilter].
func (this *SEntry) SetPathFilter(filter FPathFilter) {
	this.pathFilter = filter
	this.pathFilterCache = make([]string, 0)
}

// Updates the info and infoError for just this entry.
// Does not handle any infoError that may be returned.
// Does not Update any children.
//
// In case that the new info returns nil,
// it will not override the existing info.
func (this *SEntry) Update() {
	var info, err = this.getInfo()

	this.infoError = err

	if info != nil {
		this.info = info
	}
}

func (this *SEntry) SortChildren() {
	slices.SortStableFunc(this.children, func(a *SEntry, b *SEntry) int {
		var prio = int(a.info.Mode()) - int(b.info.Mode())
		return prio
	})
}

func (this *SEntry) removeChild(index int) {
	this.children = slices.Delete(this.children, index, index+1)
	// this.children = append(this.children[:index], this.children[index+1:]...)
}

// Runs fn for each child in this and removes them if fn returns false.
func (this *SEntry) CleanChildListFunc(fn func(child *SEntry) bool) {
	for i := 0; i < len(this.children); i++ {
		var child = this.children[i]

		if fn(child) == false {
			this.removeChild(i)
			i -= 1
		}
	}
}

// Applies the current pathFilter on all children
// and removes the children that do not meet it.
func (this *SEntry) ApplyPathFilter() {
	this.CleanChildListFunc(func(child *SEntry) bool {
		var childPath = child.GetPath()

		if this.pathFilter == nil || this.pathFilter(childPath, child.info) {
			return true
		}

		this.pathFilterCache = append(this.pathFilterCache, childPath)

		return false
	})
}

// Looks for entries in this directory that do not yet
// exist in children and adds them as a entry to children.
//
// For each new entry created, UpdateChildList will be called on that entry.
//
// If an error occurs while reading this directory it logs the error
// after processing the entries that did return before the error.
func (this *SEntry) UpdateChildList() {
	if !this.IsDirectory() {
		return
	}

	var subEntries, err = os.ReadDir(this.GetPath())

	for _, subEntry := range subEntries {
		if this.GetChildByName(subEntry.Name()) != nil {
			continue
		}

		var newPath = filepath.Join(this.GetPath(), subEntry.Name())
		var entryInfo, err = subEntry.Info()

		if err != nil {
			clog.Errorf(
				"Unexpected error while updating '%s'\n'%s' will be skipped\nError: %v",
				this.GetPath(),
				subEntry.Name(),
				err,
			)
			continue
		}

		if this.pathFilter != nil &&
			(slices.Contains(this.pathFilterCache, newPath) ||
				this.pathFilter(newPath, entryInfo) == false) {
			this.pathFilterCache = append(this.pathFilterCache, newPath)
			continue
		}

		var entry = newEntry(newPath)

		entry.SetPathFilter(this.pathFilter)
		entry.UpdateChildList()

		this.children = append(this.children, entry)
	}

	this.SortChildren()

	if err != nil {
		clog.Errorf("Error while reading directory %s:\n%v", this.GetPath(), err)
		clog.Debugf("Entries returned before error:\n%v", subEntries)
	}
}

func (this *SEntry) ForEachChild(fn func(child *SEntry)) {
	for _, child := range this.children {
		fn(child)
	}
}

// Calls fn for each child in children.
func (this *SEntry) Recursive(fn func(entry *SEntry), includeSelf bool) {
	if includeSelf {
		fn(this)
	}

	this.ForEachChild(func(child *SEntry) {
		fn(child)
		child.Recursive(fn, false)
	})
}

func (this SEntry) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		this.info.Mode(),
		time.Since(this.info.ModTime()).Round(time.Second),
		this.info.Name(),
	)
}

func (this SEntry) StringRecursive() string {
	var modeLines = []string{this.info.Mode().String()}
	var timeLines = []string{this.info.ModTime().Round(time.Second).String()}
	var nameLines = []string{}

	var modeLineSize = len(modeLines[0])
	var timeLineSize = len(timeLines[0])

	const gap = 2

	var getName = func(entry SEntry) string {
		var builder strings.Builder

		var name = entry.GetPath()
		var rel, err = filepath.Rel(this.GetPath(), name)

		if err != nil {
			builder.WriteString(name)
		} else {
			builder.WriteString(rel)
		}

		if entry.IsDirectory() {
			builder.WriteRune('/')
		}

		return builder.String()
	}

	this.Recursive(func(child *SEntry) {
		var mode = child.info.Mode().String()
		var time = child.info.ModTime().Round(time.Second).String()

		if len(mode) > modeLineSize {
			modeLineSize = len(mode)
		}
		if len(time) > timeLineSize {
			timeLineSize = len(time)
		}

		modeLines = append(modeLines, mode)
		timeLines = append(timeLines, time)
		nameLines = append(nameLines, getName(*child))
	}, true)

	var builder strings.Builder

	for i, name := range nameLines {
		var mode = modeLines[i]
		var time = timeLines[i]

		if builder.Len() > 0 {
			builder.WriteRune('\n')
		}

		builder.WriteString(mode)
		builder.WriteString(strings.Repeat(" ", modeLineSize-len(mode)+gap))
		builder.WriteString(time)
		builder.WriteString(strings.Repeat(" ", timeLineSize-len(time)+gap))
		builder.WriteString(name)
	}

	return builder.String()
}
