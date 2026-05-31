package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
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
		clog.Debug("got all files:\n")
		for _, item := range *result {
			clog.Messagef("%s:\n%s", item.Filename, ctnstring.Indent(item.Content, 2, " "))
		}
	})

	// cli.CommandWatcher()
}
