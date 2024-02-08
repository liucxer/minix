package pkg

// DataBlockTable 数据块表 总计占用空间65535K大小
type DataBlockTable struct {
	InodeItems []Inode `json:"inodeItems"`
}
