package config

import (
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

type SConfigHandlersNSDefinitions struct {
	Enable      bool
	Destination string
}

func (this *SConfigHandlersNSDefinitions) ValidateValues() error {
	var path, err = ctnfile.ValidateFilePath(filepath.Dir(constants.ConfigFilePath), this.Destination)
	this.Destination = path
	return err
}

type SConfigHandlers struct {
	NSDefinitions *SConfigHandlersNSDefinitions
}

func (this *SConfigHandlers) ValidateValues() error {
	var fields = []IConfigGroup{
		this.NSDefinitions,
	}

	for _, field := range fields {
		var err = field.ValidateValues()
		if err != nil {
			clog.Errorf("Unable to validate field: %T", field)
			return err
		}
	}

	return nil
}
