package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
	ctnstruct "github.com/CTNOriginals/CTNGoUtils/v2/struct"
)

type TConfigFilrPatterns struct {
	Include []string
	Exclude []string
}

type TConfig struct {
	Port             string
	Directory        string
	FileScanInterval int
	FilePatterns     TConfigFilrPatterns
}

var Values = &TConfig{
	Port:             "8080",
	Directory:        "./",
	FileScanInterval: 100,
	FilePatterns: TConfigFilrPatterns{
		Include: []string{"*.js", "*.ts"},
		Exclude: []string{"*.d.ts"},
	},
}

func Initialize() {
	var err error = nil

	var content []byte
	if content, err = toml.Marshal(Values); err != nil {
		panic(fmt.Sprintf("Default config values marshal error:\n%v\n", err))
	}

	log(fmt.Sprintf("Defaults:\n%s\n", content))

	if !ctnfile.FileExists(constants.ConfigFilePath) {
		var content, _ = toml.Marshal(Values)
		ctnfile.WriteFile(constants.ConfigFilePath, strings.Split(string(content), "\n"))
	}

	var meta toml.MetaData
	if meta, err = toml.DecodeFile(constants.ConfigFilePath, &Values); err != nil {
		fmt.Printf("Metadata:\n%s\n\n", ctnstruct.ToString(meta))
		panic(fmt.Sprintf("Config decode error:\n%v\n", err))
	}

	log(fmt.Sprintf("Config file content:\n%+v\n", Values))
	validateConfigValues()
	log(fmt.Sprintf("Config Values:\n%+v\n", Values))
}

func validateConfigValues() {
	ValidateBitburnerDirectory(Values.Directory)
}

func log(msg string) {
	if !constants.Debug {
		return
	}

	print(msg)
}
