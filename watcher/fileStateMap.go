package watcher

import (
	"os"
	"path/filepath"
	"time"
)

type MFileState map[string]time.Time

func (this *MFileState) Push(path string) error {
	var info, err = os.Stat(path)

	if err != nil {
		return err
	}

	(*this)[path] = info.ModTime()

	return nil
}

func (this MFileState) GetNewEntries(dir string, fn func(path string)) {
	var entries, err = os.ReadDir(dir)
	var subDirs = make([]string, 0)

	for _, entry := range entries {
		var path = filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			subDirs = append(subDirs, path)
			continue
		}

		var _, exists = this[path]

		if exists || !filePatternFilter(path) {
			continue
		}

		fn(path)
	}

	if err != nil {
		clog.Errorf("Encountered and error while reading directory %s:\n%v", dir, err)
		clog.Debugf("Entries returned before error: %v", entries)
	}

	for _, subDir := range subDirs {
		this.GetNewEntries(subDir, fn)
	}
}
