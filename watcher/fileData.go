package watcher

import (
	"os"
	"path/filepath"
)

// A wrapper for [os.FileInfo] to add some fields needed tot track a files state.
type SFileData struct {
	os.FileInfo

	dir string
}

func GetFileData(path string) (SFileData, error) {
	var info, err = os.Stat(path)

	var data = SFileData{
		FileInfo: info,
		dir:      filepath.Dir(path),
	}

	return data, err
}
