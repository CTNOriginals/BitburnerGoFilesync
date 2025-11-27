package definitions

type Definitions map[Method]Definition

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
	GetDefinitionFile: {
		Method:     GetDefinitionFile,
		Parameters: ParameterFields{},
		Response:   Response{Typ: ResString},
	},

	GetAllServers: {
		Method:     GetAllServers,
		Parameters: ParameterFields{},
		Response: Response{
			Typ:     ResObject,
			IsArray: true,
			Fields: []ResponseField{
				{Field: "hostname", Typ: ResString},
				{Field: "hasAdminRights", Typ: ResBoolean},
				{Field: "purchasedByPlayer", Typ: ResBoolean},
			},
		},
	},
}
