package watcher

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")
	var node = newNode(config.Values.Directory)

	if node.infoError != nil {
		clog.Error(node.infoError)
	}

	clog.Infof("Node:\n%s", ctnstruct.ToString(node))

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
