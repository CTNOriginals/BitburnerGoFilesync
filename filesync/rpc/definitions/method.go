package definitions

type (
	Method string
)

const (
	MethodError        Method = "error"
	PushFile           Method = "pushFile"
	GetFile            Method = "getFile"
	GetFileMetadata    Method = "getFileMetadata"
	DeleteFile         Method = "deleteFile"
	GetFileNames       Method = "getFileNames"
	GetAllFiles        Method = "getAllFiles"
	GetAllFileMetadata Method = "getAllFileMetadata"
	CalculateRam       Method = "calculateRam"
	GetDefinitionFile  Method = "getDefinitionFile"
	GetSaveFile        Method = "getSaveFile"
	GetAllServers      Method = "getAllServers"
)

func MethodsAsArray() []Method {
	return []Method{
		MethodError,
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
	}
}
