package shared

import (
	"os"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

type SetId struct {
	Host string
}

type BoolReply struct {
	Status bool
}

type FileInfoArg struct {
	SetId	SetId
	FileInfo FileInfo
}

type FileInfoListReply struct {
	Files map[string]FileInfo
}

type FileInitArg struct {
	SetId	SetId
	FileInfo FileInfo
}

type FileInitReply struct {
	Id int
}

type FileBlockArg struct {
	Id int
	Size int
	Data []byte
}

type FileCloseArg struct {
	Id int
}

const (
	//BLOCK_SIZE = 512 * 1024
	BLOCK_SIZE = 1024
)