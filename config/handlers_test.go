package config

import (
	"path/filepath"
	"testing"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestHandlersNSDefinitionsValidateValues_Valid(t *testing.T) {
	var savedConfigPath = constants.ConfigFilePath
	t.Cleanup(func() { constants.ConfigFilePath = savedConfigPath })

	var dir = t.TempDir()
	constants.ConfigFilePath = filepath.Join(dir, "config.toml")

	var nsd = SConfigHandlersNSDefinitions{
		Destination: dir,
	}

	var err = nsd.ValidateValues()
	if err != nil {
		t.Errorf("ValidateValues() = %v, want nil", err)
	}
	if nsd.Destination != dir {
		t.Errorf("Destination = %q, want %q", nsd.Destination, dir)
	}
}

func TestHandlersNSDefinitionsValidateValues_Invalid(t *testing.T) {
	var savedConfigPath = constants.ConfigFilePath
	t.Cleanup(func() { constants.ConfigFilePath = savedConfigPath })

	var dir = t.TempDir()
	constants.ConfigFilePath = filepath.Join(dir, "config.toml")

	var nsd = SConfigHandlersNSDefinitions{
		Destination: "/nonexistent/path/that/does/not/exist",
	}

	var err = nsd.ValidateValues()
	if err == nil {
		t.Error("ValidateValues() = nil, want error for nonexistent path")
	}
}

func TestHandlersValidateValues(t *testing.T) {
	var savedConfigPath = constants.ConfigFilePath
	t.Cleanup(func() { constants.ConfigFilePath = savedConfigPath })

	var dir = t.TempDir()
	constants.ConfigFilePath = filepath.Join(dir, "config.toml")

	var handlers = &SConfigHandlers{
		NSDefinitions: &SConfigHandlersNSDefinitions{
			Destination: "/nonexistent/path",
		},
	}

	var err = handlers.ValidateValues()
	if err == nil {
		t.Error("ValidateValues() = nil, want error propagated from NSDefinitions")
	}
}
