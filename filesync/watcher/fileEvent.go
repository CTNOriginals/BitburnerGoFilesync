package watcher

type FileEvent int

const (
	OnFileModify FileEvent = iota
	OnFileCreate
	OnFileDelete
)

type FileEventHandler map[FileEvent]func(file *FileInfo)

func (this FileEventHandler) Handle(file *FileInfo, event FileEvent) {
	fn := this[event]
	fn(file)
}
