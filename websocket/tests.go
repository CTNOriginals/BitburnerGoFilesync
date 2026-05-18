package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket/rpcschema"
)

func TestCommunication() {
	go StartServer(config.Values.Port)

	for {
		if ActiveConnection != nil {
			break
		}
	}

	SendRequest(rpcschema.GetAllFiles, func(message *Message) {
		// log.Printf("\nGetAllFiles OnResponse: {\n%s\n}\n", ctnstring.Indent(message.String(), 2, " "))
	}, "home")

	cli.CommandWatcher()
}
