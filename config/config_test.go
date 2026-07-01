package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestValidateConfigValues_Directory(t *testing.T) {
	type TCase struct {
		name      string
		directory string
		setup     func(base string)
		wantDir   string // empty signals Goexit expected
	}

	var base = t.TempDir()

	var cases = []TCase{
		{
			name:      "relative current",
			directory: "./",
			wantDir:   base,
		},
		{
			name:      "relative subdirectory",
			directory: "subdir",
			setup: func(base string) {
				os.MkdirAll(filepath.Join(base, "subdir"), 0755)
			},
			wantDir: filepath.Join(base, "subdir"),
		},
		{
			name:      "absolute",
			directory: base,
			wantDir:   base,
		},
		{
			name:      "non-existent",
			directory: "/nonexistent/deadbeef",
			wantDir:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var savedValues = *Values
			var savedWorkDir = constants.WorkindDirectory
			t.Cleanup(func() {
				*Values = savedValues
				constants.WorkindDirectory = savedWorkDir
			})

			if tc.setup != nil {
				tc.setup(base)
			}

			constants.WorkindDirectory = base
			Values.Directory = tc.directory
			Values.FilePatterns = SConfigFilePatterns{
				Include: []string{"*.ts"},
			}
			Values.Handlers = &SConfigHandlers{
				NSDefinitions: &SConfigHandlersNSDefinitions{
					Enable:      true,
					Destination: base,
				},
			}

			if tc.wantDir == "" {
				// NOTE: validateConfigValues() calls runtime.Goexit() on failure.
				// Once constants has a swappable exit function (e.g. constants.Fatal),
				// replace this goroutine+timeout with a direct call that returns an error.
				var done = make(chan struct{})
				go func() {
					validateConfigValues()
					close(done)
				}()
				select {
				case <-done:
					t.Error("expected runtime.Goexit(), function returned normally")
				case <-time.After(time.Second):
				}
				return
			}

			validateConfigValues()
			if Values.Directory != tc.wantDir {
				t.Errorf("Directory = %q, want %q", Values.Directory, tc.wantDir)
			}
		})
	}
}
