package rpcschema

type Schema struct {
	Method     Method
	Parameters ParameterFields
	Response   Response
}

func (this Schema) IsError() bool {
	return this.Method == MethodError
}
