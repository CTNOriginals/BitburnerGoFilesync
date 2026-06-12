package watcher

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

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
