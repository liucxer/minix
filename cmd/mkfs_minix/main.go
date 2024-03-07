package main

import (
	"fmt"
	"github.com/liucxer/minix/internal/mkfs_minix"
	"os"

	"github.com/sirupsen/logrus"
)

func usage() {
	fmt.Println("mkfs.minix: Usage: mkfs.minix /dev/name")
	os.Exit(0)
}

const (
	blockNum = 0xffff
	magicNum = 0x138F
	nameLen  = 30 // 文件名长度, 默认30
	dirSize  = 32 // 目录大小, 默认32

	MINIX_BLOCK_SIZE_BITS = 10
	MINIX_BLOCK_SIZE      = 1 << MINIX_BLOCK_SIZE_BITS
	BITS_PER_BLOCK        = MINIX_BLOCK_SIZE << 3
	MINIX_VALID_FS        = 0x0001 /* Clean fs. */
	MINIX_ERROR_FS        = 0x0002 /* fs has errors. */

)

var (
	deviceName             = ""
	inodeNum               = uint(0)
	sectorSize             = 512 // 扇区大小
	MINIX_INODES_PER_BLOCK = uint(0)

	inodeMap = []byte{}
	zoneMap  = []byte{}
)

func SetupTables() {

}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	deviceName = os.Args[1]

	// 检查是否挂载
	isMounted, err := mkfs_minix.CheckMount(deviceName)
	if err != nil {
		return
	}
	if isMounted {
		logrus.Errorf("%s is mounted; will not make a filesystem here!", deviceName)
		return
	}

	// 初始化root block
	rootBlock := mkfs_minix.InitRootBlock()

	// 初始化boot block, todo 实现磁盘写入
	bootBlock := mkfs_minix.InitBootBlock()

	// 初始化super block, todo 实现磁盘写入
	superBlock := mkfs_minix.InitSuperBlock()

	// inodeBitmap
	// dataBlockBitmap
	// inodeTable: 根节点inode信息
	// dataBlock:  根节点数据块

	SetupTables()
	fileInfo, err := os.Stat(deviceName)
	if err != nil {
		logrus.Errorf("os.Stat err:%v, path:%s", err, deviceName)
		return
	}

	file, err := os.OpenFile(deviceName, os.O_EXCL|os.O_RDWR, os.ModePerm)
	if err != nil {
		logrus.Errorf("os.OpenFile err:%v, path:%s", err, deviceName)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	//setup_tables(); // 初始化super_block_buffer、inode_map、zone_map
	//if (check)
	//	check_blocks(); // 检查是否存在坏块，如果存在则打印出来, 检查方式: 执行lseek和read函数，看是否有报错
	//else if (listfile)
	//	get_list_blocks(listfile); // 从指定文件中读取坏块列表， 此文件每行代表一个坏块编号， 且坏块数量会被打印出来
	//
	//make_root_inode();  // 将root_block写入磁盘中inode表中, 在inode_map里面 标记已经使用的inode, 只有root node,

	//mark_good_blocks(); // 在zone_map里面 标记已经使用的block, 只有root block
	//write_tables();     // 将boot_block_buffer、super_block_buffer、inode_map、zone_map、inode_buffer写入磁盘

	_ = file
	_ = fileInfo
	_ = bootBlock
	_ = rootBlock
	_ = superBlock

	fmt.Println("11")

}
