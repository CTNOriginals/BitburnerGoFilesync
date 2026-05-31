package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	testPatterns()
	testFileStateMap()
	simulateEvents()
}

func testPatterns() {
	clog.Message("\n-- Pattern Tests --\n")

	config.Values.FilePatterns.Include = []string{"*.js", "*.ts"}
	config.Values.FilePatterns.Exclude = []string{"*.d.ts"}

	var tests = []string{
		"index.js",
		"index.ts",
		"index.d.ts",
		"foo/bar.js",
		"foo/bar.d.ts",
	}

	for _, path := range tests {
		var result = shouldIncludeFile(path)
		clog.Infof("  %-30s -> %t\n", path, result)
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

	websocket.Client.Socket.Open()
	go func() {
		for range websocket.Client.Socket.Channel {
		}
	}()

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
		websocket.Client.Socket.Close()
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
