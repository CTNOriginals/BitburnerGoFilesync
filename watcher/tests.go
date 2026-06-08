package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

func TestWatcher() {
	clog.Info("\n-- Watcher Tests --\n")

	go websocket.Client.Start(config.Values.Port)
	<-*websocket.Client.OnReadySub()

	// var dir = "/home/ctn/code/bitburner/testdir"
	// var file = "foo.js"
	// var info, _ = os.Stat(filepath.Join(dir, file))
	// var dirarr, _ = os.ReadDir(filepath.Dir(dir))
	// var dirinfo, _ = dirarr[4].Info()
	//
	// for {
	// 	clog.Debugf("File: %s\nMod: %s", info.Name(), time.Since(info.ModTime()).String())
	// 	clog.Debugf("Dir: %s\nMod: %s", dirinfo.Name(), time.Since(dirinfo.ModTime()).String())
	// 	time.Sleep(time.Second)
	// }

	// testInitialize()
	testPatterns()
	testFileStateMap()
	simulateEvents()

	for {
		time.Sleep(time.Second)
	}
}

type sInitData struct {
	Dir     string
	Include []string
	Exclude []string
}

func testInitialize() {
	var tests = []sInitData{
		{ // normal
			Dir:     config.Values.Directory,
			Include: []string{"**/*.ext"},
			Exclude: []string{"**/*.x.ext"},
		},
		{ // invalid pattern
			Dir:     config.Values.Directory,
			Include: []string{"**/*.ext\\"}, // escape without next character
			Exclude: []string{""},           // empty pattern
		},
	}

	for _, test := range tests {
		clog.Infof("Initialize test: \n%s", ctnstruct.ToString(test))
		config.Values.Directory = test.Dir
		config.Values.FilePatterns.Include = test.Include
		config.Values.FilePatterns.Exclude = test.Exclude
		Initialize()
	}
}

func testPatterns() {
	clog.Message("\n-- Pattern Tests --\n")
	// // Log all the files that the program currently sees
	// utils.ForEachFileInDirRecursive(config.Values.Directory, func(file os.FileInfo, dir string) {
	// 	clog.Debug(fmt.Sprintf("%s/%s", dir, file.Name()))
	// })

	var include = ".ext"
	var exclude = ".x.ext"

	var globs = []string{
		"*",
		"**",
		"*/*",
		"*/**",
		"**/*",
		"**/**",
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

	var paths []string
	for _, dir := range dirs {
		for _, file := range files {
			paths = append(paths, fmt.Sprintf("%s%s", dir, file))
		}
	}

	var pathWidth int
	for _, p := range paths {
		if len(p) > pathWidth {
			pathWidth = len(p)
		}
	}

	var colWidth int
	for _, glob := range globs {
		if len(glob) > colWidth {
			colWidth = len(glob)
		}
	}
	if colWidth < 5 {
		colWidth = 5
	}

	var line = fmt.Sprintf("%-*s  ", pathWidth, "Path\\Glob")
	for _, glob := range globs {
		line += fmt.Sprintf("  %-*s", colWidth, glob)
	}
	clog.Messagef("%s\n", line)

	for _, p := range paths {
		line = fmt.Sprintf("%-*s  ", pathWidth, p)
		for _, glob := range globs {
			config.Values.FilePatterns.Include = []string{glob + include}
			config.Values.FilePatterns.Exclude = []string{glob + exclude}
			Initialize()
			var state = shouldIncludeFile(p)
			var visible = fmt.Sprintf("%-*s", colWidth, fmt.Sprintf("%t", state))
			if state {
				line += fmt.Sprintf("  \x1b[34m%s\x1b[0m", visible)
			} else {
				line += fmt.Sprintf("  \x1b[31m%s\x1b[0m", visible)
			}
		}
		clog.Messagef("%s\n", line)
	}
}

func testFileStateMap() {
	clog.Message("\n-- FileStateMap Tests --\n")

	config.Values.FilePatterns.Include = []string{}
	config.Values.FilePatterns.Exclude = []string{}

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
