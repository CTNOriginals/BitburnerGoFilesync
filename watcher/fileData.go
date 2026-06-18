package watcher

import (
	"os"
)

// A wrapper for [os.FileInfo] to add some fields needed tot track a files state.
type SFileData struct {
	Path string
	os.FileInfo
}

func GetFileData(path string) (*SFileData, error) {
	var info, err = os.Stat(path)

	var data = SFileData{
		Path:     path,
		FileInfo: info,
	}

	return &data, err
}

func (this SFileData) GetInfo() (os.FileInfo, error) {
	return os.Stat(this.Path)
}

func (this SFileData) Exists() bool {
	var _, err = this.GetInfo()
	return !os.IsNotExist(err)
}

func (this SFileData) IsModified() bool {
	var info, _ = this.GetInfo()
	return this.ModTime() != info.ModTime()
}

func (this *SFileData) Update() error {
	var info, err = this.GetInfo()

	if err != nil {
		return err
	}

	this.FileInfo = info

	return nil
}
