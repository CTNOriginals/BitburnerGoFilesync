package constants

import "os"

var WorkindDirectory, _ = os.Getwd()
var BitburnerRoot = WorkindDirectory + "/../bitburner"

const (
	Port = 8080
)
