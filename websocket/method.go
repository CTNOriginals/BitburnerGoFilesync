package websocket

/*
Schemas:

	Input: {
		jsonrpc: 2.0,
		id: number,
		method: string,
		params: string | {[key: string]: string},
	}

	Output: {
		jsonrpc: 2.0,
		id: number,
		result: any,
		error: any,
	}
*/
type TMethod string

const (
	// Schema:
	//  params: filename, content, server
	//  result: OK
	PushFile TMethod = "pushFile"

	// Schema:
	//  params: filename, server
	//  result: string
	GetFile TMethod = "getFile"

	// Schema:
	//  params: filename, server
	//  result: {
	//   filename: string,
	//   atime: string,
	//   btime: string,
	//   mtime: string
	//  }
	GetFileMetadata TMethod = "getFileMetadata"

	// Schema:
	//  params: filename, server
	//  result: OK
	DeleteFile TMethod = "deleteFile"

	// Schema:
	//  params: server
	//  result: string[]
	GetFileNames TMethod = "getFileNames"

	// Schema:
	//  params: server
	//  result: {filename: string, content: string}[]
	GetAllFiles TMethod = "getAllFiles"

	// Schema:
	//  params: server
	//  result: {
	//   filename: string,
	//   atime: string,
	//   btime: string,
	//   mtime: string
	//  }[]
	GetAllFileMetadata TMethod = "getAllFileMetadata"

	// Schema:
	//  params: filename, server
	//  result: number
	CalculateRam TMethod = "calculateRam"

	// Schema:
	//  result: string
	GetDefinitionFile TMethod = "getDefinitionFile"

	// Schema:
	//  result: {
	//   identifier: string,
	//   binary: bool,
	//   save: string,
	//  }
	GetSaveFile TMethod = "getSaveFile"

	// Schema:
	//  result: {
	//   hostname: string,
	//   hasAdminRights: bool,
	//   purchasedByPlayer: bool,
	//  }[]
	GetAllServers TMethod = "getAllServers"

	// Schema:
	//  result: any
	//  error: any && !nil
	MethodError TMethod = "error"
)

func MethodsAsArray() []TMethod {
	return []TMethod{
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
