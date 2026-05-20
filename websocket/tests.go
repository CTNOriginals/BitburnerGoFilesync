package websocket

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestClient() {
	Client.Ready = make(chan struct{})
	go Client.Start(config.Values.Port)
	defer Client.Close()

	// block untill closed
	<-Client.Ready

	Client.Socket.GetAllFiles(Params_GetAllFiles{
		Server: "home",
	}, func(result *Result_GetAllFiles) {
		log.Printf("got all files: %v\n", result)
	})

	cli.CommandWatcher()
}
