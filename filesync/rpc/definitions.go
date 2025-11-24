package rpc

// The parameters of a request are always of type string.
// This type will only hold the fields of any params that need to be passed
type Request []string

type TResponse int

const (
	ResString TResponse = iota
	ResNumber
	ResBoolean
	ResObject
)

type ResponseField struct {
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

type RequestResponse struct {
	Request  Request
	Response Response
}

type MDefinitions map[Method]RequestResponse

var definitions = MDefinitions{
	GetFileNames: RequestResponse{
		Request: []string{"server"},
		Response: Response{
			Typ: ResString,
		},
	},
}
