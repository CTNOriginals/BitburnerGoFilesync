package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

type ConfigFilrPatterns struct {
	Include []string
	Exclude []string
}

type Config struct {
	Port             int
	Directory        string
	FileScanInterval int
	FilePatterns     ConfigFilrPatterns
}

var DefaultConfig = &Config{
	Port:             8080,
	Directory:        "./",
	FileScanInterval: 100,
	FilePatterns: ConfigFilrPatterns{
		Include: []string{"*.js", "*.ts"},
		Exclude: []string{"*.d.ts"},
	},
}

func Test() {
	var content, _ = toml.Marshal(DefaultConfig)
	fmt.Printf("%s", content)

	var filePath = fmt.Sprintf("%s/%s", constants.WorkindDirectory, "config.toml")
	if !ctnfile.FileExists(filePath) {
		var content, _ = toml.Marshal(DefaultConfig)
		ctnfile.WriteFile(filePath, strings.Split(string(content), "\n"))
	}

	toml.DecodeFile(filePath, &DefaultConfig)
	fmt.Printf("%+v", DefaultConfig)
}

// type ConfigValue struct {
// 	Val  any
// 	Desc []string
// }
//
// func (this ConfigValue) MarshalTOML() ([]byte, error) {
// 	// var val = this.Val
//
// 	// switch v := this.Val.(type) {
// 	// case ConfigMap:
// 	// 	// var b []byte
// 	// 	var vm = map[string]any{}
// 	// 	for key, val := range v {
// 	// 		var vb, _ = val.MarshalTOML()
// 	// 		vm[key] = vb
// 	// 	}
// 	// 	this.Val = vm
// 	// 	// return b, nil
// 	// }
//
// 	return toml.Marshal(this.Val)
// }
//
// type ConfigMap map[string]ConfigValue
//
// // func (this ConfigMap) MarshalTOML() ([]byte, error) {
// // 	var fake = make(map[string]any)
// //
// // 	for key, val := range this {
// // 		// var content, _ = val.MarshalTOML()
// // 		fake[key] = val
// // 	}
// //
// // 	var m, _ = toml.Marshal(fake)
// // 	fmt.Printf("\n%s\n", m)
// //
// // 	var un any
// // 	toml.Unmarshal(m, &un)
// // 	fmt.Printf("%v\n", un)
// //
// // 	println("")
// //
// // 	// return toml.Marshal(fake)
// // 	return toml.Marshal(un)
// // }
//
// func (this ConfigMap) Format() []byte {
//
// }
//
// // var Config = ConfigMap{
// var Config = ConfigMap{
// 	"Include": ConfigValue{
// 		Val:  []string{"*.js", "*.ts"},
// 		Desc: []string{"Include file patterns"},
// 	},
// 	"FilePatterns": ConfigValue{
// 		Val: ConfigMap{
// 			"proto": ConfigValue{
// 				Val:  "wah",
// 				Desc: []string{"waaaah!"},
// 			},
// 			"Include": ConfigValue{
// 				Val:  []string{"*.js", "*.ts"},
// 				Desc: []string{"Include file patterns"},
// 			},
// 			"Exclude": ConfigValue{
// 				Val:  []string{"*.d.ts"},
// 				Desc: []string{"Exclude file patterns"},
// 			},
// 		},
// 	},
// }
//
// func Test() {
// 	var marshel []byte = Config.Format()
// 	// marshel, _ = toml.Marshal(Config)
//
// 	fmt.Printf("%s", string(marshel))
// }
