package watcher

import (
	"os"
)

type SNode struct {
	name      string
	info      os.FileInfo
	infoError error

	children []*SNode
}

// Create a new node.
//
// name is the file path to this node.
func newNode(path string) *SNode {
	var state, err = os.Stat(path)

	var node = SNode{
		name:      path,
		info:      state,
		infoError: err,
		children:  make([]*SNode, 0),
	}

	return &node
}

func (this SNode) getInfo() (os.FileInfo, *os.PathError) {
	var info, err = os.Stat(this.name)

	if err != nil {
		var pathError, isPathError = err.(*os.PathError)

		if !isPathError {
			clog.Fatalf("Reveived unknown error:\n%v", err)
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
		if child.name == name {
			return child
		}
	}

	return nil
}

// Updates the info and infoError for just this node.
// Does not handle any infoError that may be returned.
func (this *SNode) Update() {
	this.info, this.infoError = this.getInfo()
	this.name = this.info.Name()
}

// Clean removes all children that no longer exist.
func (this *SNode) Clean() {
	for i := 0; i < len(this.children); i++ {
		var child = this.children[i]

		if !child.Exists() {
			this.children = append(this.children[:i], this.children[i+1:]...)
		}
	}
}

func (this *SNode) addChild(entry os.DirEntry) {
	var info, err = entry.Info()

	var node = SNode{
		name:      entry.Name(),
		info:      info,
		infoError: err,
	}

	this.children = append(this.children, &node)
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

	var entries, err = os.ReadDir(this.name)

	for _, entry := range entries {
		var child = this.GetChildByName(entry.Name())

		if child != nil {
			continue
		}

		this.addChild(entry)
	}

	if err != nil {
		clog.Errorf("Error while reading directory %s:\n%v", this.name, err)
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
