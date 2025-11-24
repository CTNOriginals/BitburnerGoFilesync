package definitions

type Definition struct {
	Method     Method
	Parameters ParameterFields
	Response   Response
}

func (this Definition) IsError() bool {
	return this.Method == MethodError
}

type Definitions []Definition

func (this Definitions) GetByMethod(method Method) Definition {
	for _, def := range this {
		if def.Method == method {
			return def
		}
	}

	return Definition{
		Method: MethodError,
	}
}

var RPCDefinitions = Definitions{
	{
		Method:     GetFile,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResString},
	},
	{
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

	{
		Method:     GetFileNames,
		Parameters: ParameterFields{"server"},
		Response: Response{
			Typ:     ResString,
			IsArray: true,
		},
	},

	{
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
	{
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

	{
		Method:     PushFile,
		Parameters: ParameterFields{"filename", "content", "server"},
		Response:   Response{Typ: ResString},
	},
	{
		Method:     DeleteFile,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResString},
	},

	{
		Method:     CalculateRam,
		Parameters: ParameterFields{"filename", "server"},
		Response:   Response{Typ: ResNumber},
	},

	{
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
