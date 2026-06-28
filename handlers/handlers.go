// This package handles some functionality that does
// not fit inside its own dedicated package.
package handlers

var NetscriptDefinitions *SNetscriptDefinitions

func Initialize() {
	NetscriptDefinitions = newNetscriptDefinitions()
}
