package rpc

type (
	Method string
	Param  [][2]string
)

const (
	GetFileNames Method = "getFileNames"
)
