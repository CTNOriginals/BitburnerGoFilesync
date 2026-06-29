package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

type SNSDefinitions struct{}

func newNetscriptDefinitions() *SNSDefinitions {
	var nsdef = &SNSDefinitions{}

	if config.Values.Handlers.NSDefinitions.GetOnConnect == true {
		websocket.Client.OnReadySubCallback(nsdef.clientOnConnect)
	}

	return nsdef
}

func (this SNSDefinitions) clientOnConnect() (unsub bool) {
	var dest = config.Values.Handlers.NSDefinitions.Destination

	var stat, err = os.Stat(dest)

	if err != nil {
		clog.Errorf("Unable to download NetscriptDefinitions to %s:\n%v", dest, err)
		return true
	}

	var destFile = dest

	if stat.IsDir() {
		destFile = filepath.Join(dest, "NetscriptDefinitions.d.ts")
	}

	clog.Debugf("Downloading NetscriptDefinitions to %s", dest)

	websocket.Client.Socket.GetDefinitionFile(func(result *websocket.Result_GetDefinitionFile) {
		var content = string(*result)
		ctnfile.WriteFile(destFile, strings.Split(content, "\n"))
	})

	return true
}
