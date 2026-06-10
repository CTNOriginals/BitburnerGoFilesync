package watcher

import (
	"os"
	"path/filepath"
	"runtime"
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

func (this SNode) getInfo() (os.FileInfo, *os.PathError) {
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
	return os.IsNotExist(this.infoError)
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
func (this *SNode) Update() {
	this.info, this.infoError = this.getInfo()
	// this.name = this.info.Name()
}

// CleanChildList removes all children that no longer exist.
func (this *SNode) CleanChildList() {
	for i := 0; i < len(this.children); i++ {
		var child = this.children[i]

		if !child.Exists() {
			this.children = append(this.children[:i], this.children[i+1:]...)
		}
	}
}

// Looks for entries in this directory that do not yet
// exist in children and adds them as a node to children.
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

	if err != nil {
		clog.Errorf("Error while reading directory %s:\n%v", this.GetPath(), err)
		clog.Debugf("Entries returned before error:\n%v", entries)
	}
}

// Calls fn for each child in children.
// Does not call fn for itself.
func (this *SNode) Recursive(fn func(*SNode)) {
	for _, child := range this.children {
		fn(child)
	}
}
