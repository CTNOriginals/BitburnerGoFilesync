package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	go websocket.Client.Start(config.Values.Port)
	<-*websocket.Client.OnReadySub()

	clog.Debug("OnReady!")

	testPatterns()
	testFileStateMap()
	simulateEvents()

	for {
		time.Sleep(time.Second)
	}
}

func testPatterns() {
	clog.Message("\n-- Pattern Tests --\n")

	var include = ".ext"
	var exclude = ".x.ext"

	var globs = []string{
		"*",
		"**",
		"*/*",
		"*/**",
		"**/*",
	}
	var files = []string{
		"file.ext",
		"file.x.ext",
	}
	var dirs = []string{
		"",
		"foo/",
		"foo/bar/",
	}

	for _, glob := range globs {
		config.Values.FilePatterns.Include = []string{glob + include}
		config.Values.FilePatterns.Exclude = []string{glob + exclude}

		var results = make([]string, 0)

		for _, dir := range dirs {
			for _, file := range files {
				var filepath = fmt.Sprintf("%s%s", dir, file)
				var result = shouldIncludeFile(filepath)
				results = append(results, fmt.Sprintf("%s: %t", filepath, result))
			}
		}

		clog.Infof(
			"[%s] [%s]\n%s",
			strings.Join(config.Values.FilePatterns.Include, ", "),
			strings.Join(config.Values.FilePatterns.Exclude, ", "),
			strings.Join(results, "\n"),
		)
	}

}

func testFileStateMap() {
	clog.Message("\n-- FileStateMap Tests --\n")

	Initialize()

	var lines = make([]string, 0, len(FileStateMap))
	for path := range FileStateMap {
		lines = append(lines, path)
	}

	if len(lines) > 0 {
		clog.Messagef("Registered files:\n%s\n", strings.Join(lines, "\n"))
	} else {
		clog.Message("No files registered (directory may be empty or not configured)\n")
	}

	clog.Debugf("Total files registered: %d\n", len(FileStateMap))

	if len(FileStateMap) > 0 {
		for path, state := range FileStateMap {
			clog.Messagef("  %s: modified %s\n", path, state.Info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}
}

func simulateEvents() {
	clog.Message("\n-- Simulate File Events --\n")

	var tempDir, err = os.MkdirTemp("", "watcher-test-*")
	if err != nil {
		clog.Errorf("Failed to create temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "base.js"), []byte("base"), 0644)
	os.WriteFile(filepath.Join(tempDir, "base.ts"), []byte("base"), 0644)

	var origDir = config.Values.Directory
	var origState = FileStateMap
	var origInclude = config.Values.FilePatterns.Include
	var origExclude = config.Values.FilePatterns.Exclude

	config.Values.Directory = tempDir
	config.Values.FilePatterns.Include = []string{"*"}
	config.Values.FilePatterns.Exclude = []string{}
	FileStateMap = MFileState{}

	defer func() {
		config.Values.Directory = origDir
		config.Values.FilePatterns.Include = origInclude
		config.Values.FilePatterns.Exclude = origExclude
		FileStateMap = origState
	}()

	Initialize()
	clog.Messagef("Registered: %d files\n", len(FileStateMap))

	os.WriteFile(filepath.Join(tempDir, "new.js"), []byte("new file"), 0644)
	scanFiles()

	time.Sleep(100 * time.Millisecond)
	os.WriteFile(filepath.Join(tempDir, "base.js"), []byte("modified content"), 0644)
	scanFiles()

	os.Remove(filepath.Join(tempDir, "base.ts"))
	scanFiles()

	clog.Messagef("Final state: %d files\n", len(FileStateMap))
}
