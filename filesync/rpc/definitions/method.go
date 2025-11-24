package definitions

type (
	Method string
)

const (
	MethodError  Method = "error"
	GetFileNames Method = "getFileNames"
)

func MethodsAsArray() []Method {
	return []Method{
		MethodError,
		GetFileNames,
	}
}
