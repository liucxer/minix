package mkfs_minix

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/kent007/linux-inspect/etc"
	"github.com/liucxer/minix/internal/pkg"
	"github.com/sirupsen/logrus"
	"unsafe"
)

const (
	MINIX_BLOCK_SIZE_BITS = 10
	MINIX_BLOCK_SIZE      = 1 << MINIX_BLOCK_SIZE_BITS

	BITS_PER_BLOCK = MINIX_BLOCK_SIZE << 3
	MINIX_VALID_FS = 0x0001 /* Clean fs. */
	MINIX_ERROR_FS = 0x0002 /* fs has errors. */

	blockNum = 0xffff

	magicNum = 0x138F
)

var (
	MINIX_INODES_PER_BLOCK = uint(0)
	inodeNum               = uint(0)
)

func CheckMount(deviceName string) (bool, error) {
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

func InitRootBlock() [MINIX_BLOCK_SIZE]byte {
	var rootBlock = [MINIX_BLOCK_SIZE]byte{}
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
	copy(rootBlock[32*2+4:], badBlock)
	return rootBlock
}

func UPPER(size, n int) int {
	return (size + ((n) - 1)) / (n)
}

func InitSuperBlock() pkg.SuperBlock {
	nodeSize := uint(unsafe.Sizeof(pkg.Inode{}))
	MINIX_INODES_PER_BLOCK = MINIX_BLOCK_SIZE / nodeSize
	inodeNum = blockNum / 3
	inodeNum = (inodeNum + MINIX_INODES_PER_BLOCK - 1) & ^(MINIX_INODES_PER_BLOCK - 1)

	var superBlock pkg.SuperBlock
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
	return superBlock
}

func InitBootBlock() [MINIX_BLOCK_SIZE]byte {
	bootBootBlock := [MINIX_BLOCK_SIZE]byte{}

	return bootBootBlock
}

func InitInodeBitMap(superBlock pkg.SuperBlock) bitset.BitSet {
	var b bitset.BitSet

	for i := 0; i < int(superBlock.InodeBitmapBlocksNum)*MINIX_BLOCK_SIZE*8; i++ {
		b.Set(uint(i))
	}
	for i := 0; i < int(superBlock.InodeNum); i++ {
		b.Clear(uint(i))
	}
	return b
}

func InitDataBlockBitMap() {

}
