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

type SNode struct {
	// The directory path that contains this node.
	dir string

	// The last know info of this node.
	info os.FileInfo
	// If an error occured while trying to get the info,
	// it will be stored in here.
	infoError error

	children []*SNode
}

// newNode Creates a new node from path.
// It does not pre-populate the node's list of children.
//
// If an error occurs while getting the entries info,
// the error will be stored in this node so that it
// can be evaluated with the provided functions.
func newNode(path string) *SNode {
	var info, err = os.Stat(path)

	var node = SNode{
		dir:       filepath.Dir(path),
		info:      info,
		infoError: err,
		children:  make([]*SNode, 0),
	}

	return &node
}

func (this SNode) GetPath() string {
	return filepath.Join(this.dir, this.info.Name())
}

func (this SNode) getInfo() (os.FileInfo, error) {
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

func (this SNode) Exists() bool {
	return !os.IsNotExist(this.infoError)
}

func (this SNode) IsModified() bool {
	// TODO: figure out if we need to handle an error here.
	var info, _ = this.getInfo()
	return this.info.ModTime() != info.ModTime()
}

func (this SNode) IsDirectory() bool {
	return this.info.IsDir()
}

func (this SNode) GetChildByName(name string) *SNode {
	for _, child := range this.children {
		if child.info.Name() == name {
			return child
		}
	}

	return nil
}

// Updates the info and infoError for just this node.
// Does not handle any infoError that may be returned.
//
// In case that the new info returns nil,
// it will not override the existing info.
func (this *SNode) Update() {
	var info, err = this.getInfo()

	this.infoError = err

	if info != nil {
		this.info = info
	}

	// TODO: check if the name is different from what is stored
	// and fire an event for it if it is.
}

func (this *SNode) SortChildren() {
	slices.SortStableFunc(this.children, func(a *SNode, b *SNode) int {
		var prio = int(a.info.Mode()) - int(b.info.Mode())
		return prio
	})
}

// CleanChildList removes all children that no longer exist.
func (this *SNode) CleanChildList() {
	var clean = make([]*SNode, 0)

	for _, child := range this.children {
		if child.Exists() {
			clean = append(clean, child)
		}
	}

	this.children = clean
}

// Looks for entries in this directory that do not yet
// exist in children and adds them as a node to children.
//
// For each new node created, UpdateChildList will be called on that node.
//
// If an error occurs while reading this directory it logs the error
// after processing the entries that did return before the error.
func (this *SNode) UpdateChildList() {
	if !this.IsDirectory() {
		return
	}

	var entries, err = os.ReadDir(this.GetPath())

	for _, entry := range entries {
		var child = this.GetChildByName(entry.Name())

		if child != nil {
			continue
		}

		var node = newNode(filepath.Join(this.GetPath(), entry.Name()))

		node.UpdateChildList()

		this.children = append(this.children, node)
	}

	this.SortChildren()

	if err != nil {
		clog.Errorf("Error while reading directory %s:\n%v", this.GetPath(), err)
		clog.Debugf("Entries returned before error:\n%v", entries)
	}
}

func (this *SNode) ForEachChild(fn func(child *SNode)) {
	for _, child := range this.children {
		fn(child)
	}
}

// Calls fn for each child in children.
// Does not call fn for itself.
func (this *SNode) Recursive(fn func(child *SNode)) {
	this.ForEachChild(func(child *SNode) {
		fn(child)
		child.Recursive(fn)
	})
}

func (this SNode) GetTimeSinceModify() time.Duration {
	var info, _ = this.getInfo()
	return time.Since(info.ModTime())
}

func (this SNode) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		this.info.Mode(),
		time.Since(this.info.ModTime()).Round(time.Second),
		this.info.Name(),
	)
}

func (this SNode) StringRecursive() string {
	var modeLines = []string{this.info.Mode().String()}
	var timeLines = []string{this.GetTimeSinceModify().Round(time.Second).String()}
	var nameLines = []string{}

	var modeLineSize = len(modeLines[0])
	var timeLineSize = len(timeLines[0])

	const gap = 2

	var getName = func(node SNode) string {
		var builder strings.Builder

		var name = node.GetPath()
		var rel, err = filepath.Rel(this.GetPath(), name)

		if err != nil {
			builder.WriteString(name)
		} else {
			builder.WriteString(rel)
		}

		if node.IsDirectory() {
			builder.WriteRune('/')
		}

		return builder.String()
	}

	nameLines = append(nameLines, getName(this))

	this.Recursive(func(child *SNode) {
		var mode = child.info.Mode().String()
		var time = child.GetTimeSinceModify().Round(time.Second).String()

		if len(mode) > modeLineSize {
			modeLineSize = len(mode)
		}
		if len(time) > timeLineSize {
			timeLineSize = len(time)
		}

		modeLines = append(modeLines, mode)
		timeLines = append(timeLines, time)
		nameLines = append(nameLines, getName(*child))
	})

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
