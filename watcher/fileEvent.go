package watcher

type FileEvent int

const (
	OnFileModify FileEvent = iota
	OnFileCreate
	OnFileDelete
)

type MFileEventHandler map[FileEvent]func(file *FileInfo)

func (this MFileEventHandler) Handle(file *FileInfo, event FileEvent) {
	fn := this[event]
	fn(file)
}
