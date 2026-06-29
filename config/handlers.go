package config

import (
	"os"
	"path/filepath"

	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
)

type SConfigHandlersNSDefinitions struct {
	GetOnConnect bool
	Destination  string
}

func (this *SConfigHandlersNSDefinitions) ValidateValues() error {
	var fileExists = func() error {
		var _, err = os.Stat(this.Destination)
		return err
	}

	if filepath.IsAbs(this.Destination) {
		return fileExists()
	}

	var baseDir = filepath.Dir(constants.ConfigFilePath)

	this.Destination = filepath.Clean(filepath.Join(baseDir, this.Destination))

	if filepath.IsAbs(this.Destination) {
		return fileExists()
	}

	var abspath, err = filepath.Abs(this.Destination)

	if err != nil {
		return err
	}

	this.Destination = abspath

	return fileExists()
}

type SConfigHandlers struct {
	NetscriptDefinitions *SConfigHandlersNSDefinitions
}

func (this *SConfigHandlers) ValidateValues() error {
	var fields = []IConfigGroup{
		this.NetscriptDefinitions,
	}

	for _, field := range fields {
		var err = field.ValidateValues()
		if err != nil {
			return err
		}
	}

	return nil
}
