package websocket

import "testing"

func TestMethodsAsArray(t *testing.T) {
	var methods = MethodsAsArray()
	var seen = make(map[TMethod]bool)

	for _, m := range methods {
		if seen[m] {
			t.Errorf("duplicate method: %s", m)
		}
		seen[m] = true
	}

	if len(methods) != 12 {
		t.Errorf("len = %d, want 12", len(methods))
	}

	var expected = []TMethod{
		PushFile, GetFile, GetFileMetadata, DeleteFile,
		GetFileNames, GetAllFiles, GetAllFileMetadata,
		CalculateRam, GetDefinitionFile, GetSaveFile,
		GetAllServers, MethodError,
	}
	for _, m := range expected {
		if !seen[m] {
			t.Errorf("missing method: %s", m)
		}
	}
}
