package definitions

type (
	Method string
)

const (
	// Parameters:
	//   func(filename, content, server)
	PushFile Method = "pushFile"

	// Parameters:
	//   func(filename, server)
	GetFile Method = "getFile"

	// Parameters:
	//   func(filename, server)
	GetFileMetadata Method = "getFileMetadata"

	// Parameters:
	//   func(filename, server)
	DeleteFile Method = "deleteFile"

	// Parameters:
	//   func(server)
	GetFileNames Method = "getFileNames"

	// Parameters:
	//   func(filename, server)
	GetAllFiles Method = "getAllFiles"

	// Parameters:
	//   func(server)
	GetAllFileMetadata Method = "getAllFileMetadata"

	// Parameters:
	//   func(filename, server)
	CalculateRam Method = "calculateRam"

	GetDefinitionFile Method = "getDefinitionFile"
	GetSaveFile       Method = "getSaveFile"
	GetAllServers     Method = "getAllServers"
	MethodError       Method = "error"
)

func MethodsAsArray() []Method {
	return []Method{
		PushFile,
		GetFile,
		GetFileMetadata,
		DeleteFile,
		GetFileNames,
		GetAllFiles,
		GetAllFileMetadata,
		CalculateRam,
		GetDefinitionFile,
		GetSaveFile,
		GetAllServers,
		MethodError,
	}
}
