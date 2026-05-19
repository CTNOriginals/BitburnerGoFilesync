package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestCommunication() {
	go StartServer(config.Values.Port)

	for {
		if ActiveConnection != nil {
			break
		}
	}

	cli.CommandWatcher()
}
