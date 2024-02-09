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

	// 初始化boot block
	bootBlock := mkfs_minix.InitBootBlock()

	// 初始化super block
	superBlock := mkfs_minix.InitSuperBlock()

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

	_ = file
	_ = fileInfo
	_ = bootBlock
	_ = rootBlock
	_ = superBlock
	SetupTables()
	fmt.Println("11")
}
