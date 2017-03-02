package utils

import ()

// IDX_ROOT_PATH 默认索引存放位置
const IDX_ROOT_PATH string = "./index/"

const IDX_DETAIL_PATH string = "./detail/"
const (
	MODE_APPEND = iota
	MODE_CREATE
)

// SimpleFieldInfo description: 字段的描述信息
type SimpleFieldInfo struct {
	FieldName string `json:"fieldname"`
	FieldType uint64 `json:"fieldtype"`
	PflOffset uint64 `json:"pfloffset"` //正排索引的偏移量
	PflLen    int    `json:"pfllen"`    //正排索引长度
}

// IndexStrct 索引构造结构，包含字段信息
type IndexStrct struct {
	IndexName    string            `json:"indexname"`
	IndexMapping []SimpleFieldInfo `json:"indexmapping"`
}

type Mmap struct {
	MmapBytes   []byte
	FileName    string
	FileLen     int64
	FilePointer int64
	MapType     int64
	FileFd      *os.File
}

func NewMmap(file_name string, mode int) (*Mmap, error) {
	this := &Mmap{MmapBytes: make([]byte, 0), FileName: file_name, FileLen: 0, MapType: 0, FilePointer: 0, FileFd: nil}
	file_mode := os.O_RDWR
	f, err := os.OpenFile(file_name, file_mode, 0664)
}
