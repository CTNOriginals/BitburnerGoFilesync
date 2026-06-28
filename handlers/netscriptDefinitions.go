package handlers

import (
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
	// TODO: add dedicated config field for nsdef destination dir
	var dest = config.Values.Directory

	clog.Debugf("Downloading NetscriptDefinitions to %s", dest)

	websocket.Client.Socket.GetDefinitionFile(func(result *websocket.Result_GetDefinitionFile) {
		var content = string(*result)
		ctnfile.WriteFile(filepath.Join(dest, "NetscriptDefinitions.d.ts"), strings.Split(content, "\n"))
	})

	return true
}
