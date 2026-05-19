package websocket

type ServerName struct {
	Server string `json:"server"`
}
type File_Content struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}
type File_Server struct {
	Filename string `json:"filename"`
	Server   string `json:"server"`
}
type File_Content_Server struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	Server   string `json:"server"`
}
type FileMetadata struct {
	Filename string `json:"filename"`
	Atime    string `json:"atime"`
	Btime    string `json:"btime"`
	Mtime    string `json:"mtime"`
}
type SaveFile struct {
	Identifier string `json:"identifier"`
	Binary     string `json:"binary"`
	Save       string `json:"save"`
}
type AllServers struct {
	Hostname          string `json:"hostname"`
	HasAdminRights    string `json:"hasAdminRights"`
	PurchasedByPlayer string `json:"purchasedByPlayer"`
}

type Params_PushFile File_Content_Server
type Result_PushFile string

type Params_GetFile File_Server
type Result_GetFile string

type Params_GetFileMetadata File_Server
type Result_GetFileMetadata FileMetadata

type Params_DeleteFile File_Server
type Result_DeleteFile string

type Params_GetFileNames ServerName
type Result_GetFileNames []string

type Params_GetAllFiles ServerName
type Result_GetAllFiles File_Content

type Params_GetAllFileMetadata ServerName
type Result_GetAllFileMetadata []Result_GetAllFileMetadata

type Params_CalculateRam File_Server
type Result_CalculateRam float64

// type Params_GetDefinitionFile any
type Result_GetDefinitionFile string

// type Params_GetSaveFile any
type Result_GetSaveFile SaveFile

// type Params_GetAllServers any
type Result_GetAllServers AllServers
