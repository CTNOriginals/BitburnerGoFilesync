package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

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
		log.Panicf("Default config values, marshal error:\n%v\n", err)
	}

	logConfig(fmt.Sprintf("Defaults:\n%s\n", content))

	if !ctnfile.FileExists(constants.ConfigFilePath) {
		var content, _ = toml.Marshal(Values)
		ctnfile.WriteFile(constants.ConfigFilePath, strings.Split(string(content), "\n"))
	}

	if _, err = toml.DecodeFile(constants.ConfigFilePath, &Values); err != nil {
		log.Panicf("Config decode error:\n%v", err)
	}

	logConfig(fmt.Sprintf("Config file content:\n%+v\n", Values))
	validateConfigValues()
	logConfig(fmt.Sprintf("Config Values:\n%+v\n", Values))
}

func validateConfigValues() {
	ValidateBitburnerDirectory(Values.Directory)
}

func logConfig(msg string) {
	if !constants.Debug || !constants.LogConfig {
		return
	}

	log.Print(msg)
}
