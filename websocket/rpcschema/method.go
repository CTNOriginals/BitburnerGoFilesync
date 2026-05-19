package rpcschema

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
type Method string

const (
	// Schema:
	//  params: filename, content, server
	//  result: OK
	PushFile Method = "pushFile"

	// Schema:
	//  params: filename, server
	//  result: string
	GetFile Method = "getFile"

	// Schema:
	//  params: filename, server
	//  result: {
	//   filename: string,
	//   atime: string,
	//   btime: string,
	//   mtime: string
	//  }
	GetFileMetadata Method = "getFileMetadata"

	// Schema:
	//  params: filename, server
	//  result: OK
	DeleteFile Method = "deleteFile"

	// Schema:
	//  params: server
	//  result: string[]
	GetFileNames Method = "getFileNames"

	// Schema:
	//  params: server
	//  result: {filename: string, content: string}[]
	GetAllFiles Method = "getAllFiles"

	// Schema:
	//  params: server
	//  result: {
	//   filename: string,
	//   atime: string,
	//   btime: string,
	//   mtime: string
	//  }[]
	GetAllFileMetadata Method = "getAllFileMetadata"

	// Schema:
	//  params: filename, server
	//  result: number
	CalculateRam Method = "calculateRam"

	// Schema:
	//  result: string
	GetDefinitionFile Method = "getDefinitionFile"

	// Schema:
	//  result: {
	//   identifier: string,
	//   binary: bool,
	//   save: string,
	//  }
	GetSaveFile Method = "getSaveFile"

	// Schema:
	//  result: {
	//   hostname: string,
	//   hasAdminRights: bool,
	//   purchasedByPlayer: bool,
	//  }[]
	GetAllServers Method = "getAllServers"

	// Schema:
	//  result: any
	//  error: any && !nil
	MethodError Method = "error"
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
