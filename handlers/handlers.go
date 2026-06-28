// This package handles some functionality that does
// not fit inside its own dedicated package.
package handlers

import "github.com/CTNOriginals/BitburnerGoFilesync/clogger"

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "handlers",
})

var NetscriptDefinitions *SNetscriptDefinitions

func Initialize() {
	NetscriptDefinitions = newNetscriptDefinitions()
}
