package watcher

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	var node = newNode(config.Values.Directory)

	if node.infoError != nil {
		clog.Error(node.infoError)
	}

	node.UpdateChildList()
	// node.Recursive((*SNode).UpdateChildList)
	node.ForEachChild(func(child *SNode) {
		child.Recursive((*SNode).UpdateChildList)
	})

	clog.Infof("Node:\n%s", node.StringRecursive())

	// var arr = []string{"a", "b", "c", "d", "e"}

	// for i, item := range arr {
	// 	if item == "c" {
	// 		arr = append(arr[:i], arr[i+1:]...)
	// 	}
	// 	clog.Infof("%d: %s", i, item)
	// }

	// for i := 0; i < len(arr); i++ {
	// 	var item = arr[i]
	//
	// 	if item == "c" {
	// 		arr = append(arr[:i], arr[i+1:]...)
	// 		continue
	// 	}
	// 	clog.Infof("%d: %s", i, item)
	//
	// }
	// clog.Infof("%v", arr)
}
