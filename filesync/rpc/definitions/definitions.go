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
		Method:     GetFileNames,
		Parameters: ParameterFields{"server"},
		Response: Response{
			Typ:     ResString,
			IsArray: true,
		},
	},
}
