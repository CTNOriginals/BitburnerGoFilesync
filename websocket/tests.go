package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestCommunication() {
	go Client.Start(config.Values.Port)

	for {
		if Client.Active() {
			break
		}
	}

	cli.CommandWatcher()
}
