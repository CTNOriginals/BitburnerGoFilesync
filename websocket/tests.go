package websocket

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestCommunication() {
	go Client.Start(config.Values.Port)
	defer Client.Close()

	for {
		if Client.Active() {
			break
		}
	}

	Client.Socket.GetAllFiles(getAllFilesParams{
		Server: "home",
	}, func(result []getAllFilesResult) {
		log.Printf("got all files: %v\n", result)
	})

	cli.CommandWatcher()
}
