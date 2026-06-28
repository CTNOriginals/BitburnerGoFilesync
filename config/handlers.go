package config

import (
	"os"
	"path/filepath"
)

type SConfigHandlersNSDefinitions struct {
	GetOnConnect bool
	Destination  string
}

func (this *SConfigHandlersNSDefinitions) ValidateValues() error {
	var fileExists = func() error {
		var stat, err = os.Stat(this.Destination)
		clog.Debugf("nsdef path: %s", stat.Name())
		return err
	}

	if filepath.IsAbs(this.Destination) {
		return fileExists()
	}

	// TODO: finish this
	// this.Destination = filepath.Clean(this.Destination)
	//
	// var relpath, err = filepath.Rel(Values.Directory, this.Destination)
	//
	// if err != nil {
	// 	return err
	// }
	//
	// this.Destination = relpath

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
