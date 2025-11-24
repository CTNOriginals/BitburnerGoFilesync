package definitions

type TResponse int

const (
	ResString TResponse = iota
	ResNumber
	ResBoolean
	ResObject
)

type ResponseField struct {
	// The name/key of the field
	Field string
	// A field type should never be of type object
	Typ TResponse
}

type Response struct {
	Typ TResponse
	// Indicates that the result will be an array of Typ
	IsArray bool
	// This field should only be populated if Typ == Object
	Fields []ResponseField
}
