package constants

import (
	"fmt"
	"os"
)

var WorkindDirectory, _ = os.Getwd()
var ConfigFilePath = fmt.Sprintf("%s/%s", WorkindDirectory, "config.toml")

var Debug = false
var NoWatcher = false
var NoServer = false

var KeepAlive = false

var LogConfig = false
