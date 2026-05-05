package watcher

import (
	"testing"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	. "github.com/smartystreets/goconvey/convey"
)

func testFilePattern(path string, state bool) {
	var assertion = ShouldBeTrue

	if !state {
		assertion = ShouldBeFalse
	}

	Convey(path, func() {
		So(shouldIncludeFile(path), assertion)
	})
}

func TestPatternMatching(t *testing.T) {
	Convey("Config values are populated", t, func() {
		Convey("Include and Exclude are empty, all should be true", func() {
			config.Values.FilePatterns.Include = []string{}
			config.Values.FilePatterns.Exclude = []string{}

			testFilePattern("index.ts", true)
			testFilePattern("index.d.ts", true)
			testFilePattern("dir/index.ts", true)
			testFilePattern("dir/index.d.ts", true)
			testFilePattern("foo/bar/dir/index.ts", true)
			testFilePattern("foo/bar/dir/index.d.ts", true)
		})

		Convey("Include *.ts Exclude *.d.ts", func() {
			config.Values.FilePatterns.Include = []string{
				"*.ts",
			}
			config.Values.FilePatterns.Exclude = []string{
				"*.d.ts",
			}

			Convey("Should return true", func() {
				testFilePattern("index.ts", true)
				testFilePattern("dir/index.ts", true)
				testFilePattern("foo/bar/dir/index.ts", true)
			})
			Convey("Should return false", func() {
				testFilePattern("index.d.ts", false)
				testFilePattern("dir/index.d.ts", false)
				testFilePattern("foo/bar/dir/index.d.ts", false)
			})
		})

		Convey("Include **/*.ts Exclude **/*.d.ts", func() {
			config.Values.FilePatterns.Include = []string{
				"**/*.ts",
			}
			config.Values.FilePatterns.Exclude = []string{
				"**/*.d.ts",
			}

			Convey("Should return true", func() {
				testFilePattern("dir/index.ts", true)
				testFilePattern("foo/bar/dir/index.ts", true)
			})
			Convey("Should return false", func() {
				testFilePattern("index.ts", false)
				testFilePattern("index.d.ts", false)
				testFilePattern("dir/index.d.ts", false)
				testFilePattern("foo/bar/dir/index.d.ts", false)
			})
		})

		Convey("Include **/foo/** Exclude **/foo/bar/**", func() {
			config.Values.FilePatterns.Include = []string{
				"foo/**",
				"**/foo/**",
			}
			config.Values.FilePatterns.Exclude = []string{
				"foo/bar/**",
				"**/foo/bar/**",
			}

			Convey("Should return true", func() {
				testFilePattern("foo/index.ts", true)
				testFilePattern("root/foo/child/main.js", true)
				testFilePattern("root/dir/foo/index.html", true)
			})
			Convey("Should return false", func() {
				testFilePattern("index.ts", false)
				testFilePattern("root/index.d.ts", false)
				testFilePattern("foo/bar/index.ts", false)
				testFilePattern("root/foo/bar/index.ts", false)
			})

		})
	})
}
