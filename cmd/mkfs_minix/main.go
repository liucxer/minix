package main

import (
	"fmt"
	"github.com/liucxer/minix/internal/pkg"
	"os"
	"unsafe"

	"github.com/kent007/linux-inspect/etc"
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
	rootBlock              = [MINIX_BLOCK_SIZE]byte{}
	sectorSize             = 512 // 扇区大小
	MINIX_INODES_PER_BLOCK = uint(0)
	superBlock             = pkg.SuperBlock{}

	inodeMap = []byte{}
	zoneMap  = []byte{}
)

func checkMount() (bool, error) {
	var err error
	var isMounted bool
	mss, err := etc.GetMtab()
	if err != nil {
		logrus.Errorf("etc.GetMtab err:%v", err)
		return isMounted, err
	}
	for _, ms := range mss {
		if ms.FileSystem == deviceName {
			isMounted = true
		}
	}

	return isMounted, err
}

func initRootBlock() {
	// 初始化 root_block 会写入block区域
	//      V1、V2版本
	//      |--------|--------|--------|--------| dirsize=32
	//      |-1--.---|--------|--------|--------|
	//      |-1--..--|--------|--------|--------|
	//      |-2--.bad|blocks--|--------|--------|
	rootBlock[1] = 1
	rootBlock[4] = '.'

	rootBlock[32+1] = 1
	rootBlock[32+4] = '.'
	rootBlock[32+5] = '.'

	rootBlock[32*2+1] = 2
	badBlock := ".badblocks"
	copy(rootBlock[32*2+4:], []byte(badBlock))
}

func UPPER(size, n int) int {
	return (size + ((n) - 1)) / (n)
}

func initSuperBlock() {
	nodeSize := uint(unsafe.Sizeof(pkg.Inode{}))
	MINIX_INODES_PER_BLOCK = MINIX_BLOCK_SIZE / nodeSize
	inodeNum = blockNum / 3
	inodeNum = (inodeNum + MINIX_INODES_PER_BLOCK - 1) & ^(MINIX_INODES_PER_BLOCK - 1)
	superBlock.InodeNum = uint16(inodeNum)
	superBlock.ZoneNum = blockNum
	superBlock.InodeBitmapBlocksNum = uint16(UPPER(int(inodeNum+1), BITS_PER_BLOCK))
	superBlock.ZoneBitmapBlocksNum = uint16(UPPER(blockNum-(1+int(superBlock.InodeBitmapBlocksNum)+UPPER(int(superBlock.InodeNum), int(MINIX_INODES_PER_BLOCK))), BITS_PER_BLOCK+1))
	superBlock.FirstDataZone = 2 + superBlock.InodeBitmapBlocksNum + superBlock.ZoneBitmapBlocksNum + uint16(UPPER(int(superBlock.InodeNum), int(MINIX_INODES_PER_BLOCK)))
	superBlock.ZoneSize = 0
	superBlock.MaxFileSize = (7 + 512 + 512*512) * 1024
	superBlock.Magic = magicNum
	superBlock.State |= MINIX_VALID_FS
	superBlock.State &= ^uint16(MINIX_ERROR_FS)
}

func setupTables() {
	initSuperBlock()

}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	deviceName = os.Args[1]

	isMounted, err := checkMount()
	if err != nil {
		return
	}
	if isMounted {
		logrus.Errorf("%s is mounted; will not make a filesystem here!", deviceName)
		return
	}

	initRootBlock()

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
	setupTables()
	fmt.Println("11")
}
