package communication

import (
	"log"

	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/communication/constructor"
	"github.com/CTNOriginals/BitburnerGoFilesync/communication/definitions"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
)

func TestCommunication() {
	go StartServer(config.Values.Port)

	for {
		if ActiveConnection != nil {
			break
		}
	}

	SendRequest(definitions.GetAllFiles, func(message *constructor.Message) {
		// log.Printf("\nGetAllFiles OnResponse: {\n%s\n}\n", ctnstring.Indent(message.String(), 2, " "))
	}, "home")

	log.Println("heaslkdjef")

	cli.CommandWatcher()
}
