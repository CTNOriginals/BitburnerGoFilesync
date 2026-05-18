package definitions

type Definition struct {
	Method     Method
	Parameters ParameterFields
	Response   Response
}

func (this Definition) IsError() bool {
	return this.Method == MethodError
}
