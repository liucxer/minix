package pkg

// InodeTable inode表 总计占用空间 683K大小
type InodeTable struct {
	InodeItems []Inode `json:"inodeItems"`
}

// Inode inode结构体 占用空间32个字节
type Inode struct {
	Mode   Mode      // 文件类型和属性(rwx 位)。
	Uid    uint16    // 用户id（文件拥有者标识符）。
	Size   uint32    // 文件大小（字节数）。
	Time   uint32    // 修改时间（自1970.1.1:0 算起，秒）。
	Gid    uint8     // 组id(文件拥有者所在的组)。
	NLinks uint8     // 链接数（多少个文件目录项指向该i 节点）。
	Zone   [9]uint16 // 直接(0-6)、间接(7)或双重间接(8)逻辑块号。7k + 512k + 512 * 512K = 262663k
}

type Mode uint16

const NameLen = 30

type DirEntry struct {
	InodeNo uint16
	Name    [NameLen]byte
}

type DirEntryList []DirEntry
