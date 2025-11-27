package definitions

type Definition struct {
	Method     Method
	Parameters ParameterFields
	Response   Response
}

func (this Definition) IsError() bool {
	return this.Method == MethodError
}

type Definitions map[Method]Definition

// func (this Definitions) GetByMethod(method Method) Definition {
// 	for _, def := range this {
// 		if def.Method == method {
// 			return def
// 		}
// 	}

// 	return Definition{
// 		Method: MethodError,
// 	}
// }

var RPCDefinitions = Definitions{
	GetFile: {
		Method:     GetFile,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResString},
	},
	GetAllFiles: {
		Method:     GetAllFiles,
		Parameters: ParameterFields{"filename", "server"},
		Response: Response{
			Typ:     ResObject,
			IsArray: true,
			Fields: []ResponseField{
				{Field: "filename", Typ: ResString},
				{Field: "content", Typ: ResString},
			},
		},
	},

	GetFileNames: {
		Method:     GetFileNames,
		Parameters: ParameterFields{"server"},
		Response: Response{
			Typ:     ResString,
			IsArray: true,
		},
	},

	GetFileMetadata: {
		Method:     GetFileMetadata,
		Parameters: ParameterFields{"filename", "server"},
		Response: Response{
			Typ: ResObject,
			Fields: []ResponseField{
				{Field: "filename", Typ: ResString},
				{Field: "atime", Typ: ResString},
				{Field: "btime", Typ: ResString},
				{Field: "mtime", Typ: ResString},
			},
		},
	},
	GetAllFileMetadata: {
		Method:     GetAllFileMetadata,
		Parameters: ParameterFields{"server"},
		Response: Response{
			Typ:     ResObject,
			IsArray: true,
			Fields: []ResponseField{
				{Field: "filename", Typ: ResString},
				{Field: "atime", Typ: ResString},
				{Field: "btime", Typ: ResString},
				{Field: "mtime", Typ: ResString},
			},
		},
	},

	PushFile: {
		Method:     PushFile,
		Parameters: ParameterFields{"filename", "content", "server"},
		Response:   Response{Typ: ResString},
	},
	DeleteFile: {
		Method:     DeleteFile,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResString},
	},

	CalculateRam: {
		Method:     CalculateRam,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResNumber},
	},

	GetSaveFile: {
		Method:     GetSaveFile,
		Parameters: ParameterFields{},
		Response: Response{
			Typ: ResObject,
			Fields: []ResponseField{
				{Field: "identifier", Typ: ResString},
				{Field: "binary", Typ: ResBoolean},
				{Field: "save", Typ: ResString},
			},
		},
	},
}
