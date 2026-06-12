package watcher

import "github.com/CTNOriginals/BitburnerGoFilesync/config"

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	var node = newNode(config.Values.Directory)
	// var node = newNode(constants.WorkindDirectory)

	if node.infoError != nil {
		clog.Error(node.infoError)
	}

	node.Recursive((*SNode).UpdateChildList, true)

	clog.Infof("%s:\n%s", node.GetPath(), node.StringRecursive())
}
