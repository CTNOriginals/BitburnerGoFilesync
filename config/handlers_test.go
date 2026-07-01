package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestHandlersNSDefinitionsValidateValues(t *testing.T) {
	var savedConfigPath = constants.ConfigFilePath
	t.Cleanup(func() { constants.ConfigFilePath = savedConfigPath })

	var base = t.TempDir()
	constants.ConfigFilePath = filepath.Join(base, "config.toml")

	type TCase struct {
		name        string
		destination string
		setup       func(dir string)
		useWrapper  bool
		wantErr     bool
	}

	var cases = []TCase{
		{
			name:        "existing directory",
			destination: base,
		},
		{
			name:        "existing file",
			destination: filepath.Join(base, "def.ts"),
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "def.ts"), []byte("content"), 0644)
			},
		},
		{
			name:        "relative filename",
			destination: "def.ts",
			setup: func(dir string) {
				os.WriteFile(filepath.Join(dir, "def.ts"), []byte("content"), 0644)
			},
		},
		{
			name:        "non-existent",
			destination: "/nonexistent/deadbeef",
			wantErr:     true,
		},
		{
			name:        "handler wraps nsdefinition error",
			destination: "/nonexistent/deadbeef",
			useWrapper:  true,
			wantErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(base)
			}

			var err error
			if tc.useWrapper {
				err = (&SConfigHandlers{
					NSDefinitions: &SConfigHandlersNSDefinitions{
						Destination: tc.destination,
					},
				}).ValidateValues()
			} else {
				err = (&SConfigHandlersNSDefinitions{
					Destination: tc.destination,
				}).ValidateValues()
			}

			if tc.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
