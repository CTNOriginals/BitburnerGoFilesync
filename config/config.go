package config

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
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

type TConfigFilrPatterns struct {
	Include []string
	Exclude []string
}
type TConfigLogging struct {
	NoColor bool
}

type TConfig struct {
	Port             string
	Directory        string
	FileScanInterval int
	FilePatterns     TConfigFilrPatterns
	Logging          TConfigLogging
}

var Values = &TConfig{
	Port:             "8080",
	Directory:        "./",
	FileScanInterval: 1000,
	FilePatterns: TConfigFilrPatterns{
		Include: []string{"*.js", "*.ts"},
		Exclude: []string{"*.d.ts"},
	},
	Logging: TConfigLogging{
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
}
