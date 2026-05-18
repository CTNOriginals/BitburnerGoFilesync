package websocket

import (
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket/definitions"
)

func TestCommunication() {
	go StartServer(config.Values.Port)

	for {
		if ActiveConnection != nil {
			break
		}
	}

	SendRequest(definitions.GetAllFiles, func(message *Message) {
		// log.Printf("\nGetAllFiles OnResponse: {\n%s\n}\n", ctnstring.Indent(message.String(), 2, " "))
	}, "home")

	cli.CommandWatcher()
}
