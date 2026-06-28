package config

import (
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

func init() {
	// BUG: this happens a little too late for some logs
	// that already got printed, resulting in those logs
	// containing color regardless of the setting.
	clogger.ConfigValueNoColor = &Values.Logging.NoColor
}

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "config",

	LogLevelState: clogger.MLogLevelState{
		clogger.LogInfo | clogger.LogDebug: func() bool {
			return constants.LogConfig
		},
	},
})

type IConfigGroup interface {
	ValidateValues() error
}

type SConfigLogging struct {
	NoColor bool
}

type SConfig struct {
	Port             string
	Directory        string
	FileScanInterval int
	FilePatterns     SConfigFilrPatterns
	Logging          SConfigLogging
}

var Values = &SConfig{
	Port:             "8080",
	Directory:        "./",
	FileScanInterval: 1000,
	FilePatterns: SConfigFilrPatterns{
		Include: []string{"**/*.js", "**/*.ts"},
		Exclude: []string{"**/*.d.ts"},
	},
	Logging: SConfigLogging{
		NoColor: false,
	},
}

func Initialize() {
	var err error = nil

	var content []byte
	if content, err = toml.Marshal(Values); err != nil {
		clog.Fatalf("Default config values, marshal error:\n%v\n", err)
	}

	clog.Debugf("Defaults:\n%s\n", content)

	if !ctnfile.FileExists(constants.ConfigFilePath) {
		var content, _ = toml.Marshal(Values)
		ctnfile.WriteFile(constants.ConfigFilePath, strings.Split(string(content), "\n"))
	}

	if _, err = toml.DecodeFile(constants.ConfigFilePath, &Values); err != nil {
		clog.Fatalf("Config decode error:\n%v", err)
	}

	validateConfigValues()
	clog.Debugf("Config Values:\n%+v\n", Values)
}

func validateConfigValues() {
	ValidateBitburnerDirectory(Values.Directory)

	var validated = true

	var fields = ctnstruct.Keys(Values)
	var values = ctnstruct.Values(Values)

	for i, field := range fields {
		var value = values[i]

		switch val := value.(type) {
		case IConfigGroup:
			var err = val.ValidateValues()
			if err != nil {
				clog.Errorf("Unable to validate config field: %s\n%v", field, err)
				validated = false
			}
		}
	}

	if !validated {
		clog.Fatalf("Invalid config")
		runtime.Goexit()
	}
}
