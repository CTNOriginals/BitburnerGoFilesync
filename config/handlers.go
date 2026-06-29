package config

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

type SConfigHandlersNSDefinitions struct {
	Enable      bool
	Destination string
}

func (this *SConfigHandlersNSDefinitions) ValidateValues() error {
	var path, err = ValidateFilePath(constants.ConfigFilePath, this.Destination)
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
			return err
		}
	}

	return nil
}
