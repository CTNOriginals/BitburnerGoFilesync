package config

import (
	"testing"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

func TestValidateConfigValues_Success(t *testing.T) {
	var savedValues = *Values
	var savedWorkDir = constants.WorkindDirectory
	t.Cleanup(func() {
		*Values = savedValues
		constants.WorkindDirectory = savedWorkDir
	})

	var dir = t.TempDir()
	constants.WorkindDirectory = dir
	Values.Directory = "./"
	Values.FilePatterns = SConfigFilePatterns{
		Include: []string{"*.ts"},
	}
	Values.Handlers = &SConfigHandlers{
		NSDefinitions: &SConfigHandlersNSDefinitions{
			Enable:      true,
			Destination: dir,
		},
	}

	validateConfigValues()

	if Values.Directory != dir {
		t.Errorf("Directory = %q, want %q", Values.Directory, dir)
	}
}
