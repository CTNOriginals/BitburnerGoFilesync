package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
	ctnfile "github.com/CTNOriginals/CTNGoUtils/v2/file"
)

type SNetscriptDefinitions struct{}

func newNetscriptDefinitions() *SNetscriptDefinitions {
	var nsdef = &SNetscriptDefinitions{}

	websocket.Client.OnReadySubCallback(nsdef.clientOnConnect)

	return nsdef
}

func (this SNetscriptDefinitions) clientOnConnect() (unsub bool) {
	// TODO: add a config option to disable this function from running

	var dest = config.Values.Handlers.NetscriptDefinitions.Destination

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
