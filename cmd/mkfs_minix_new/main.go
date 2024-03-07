package main

import (
	"fmt"
	"github.com/kent007/linux-inspect/etc"
	"github.com/liucxer/minix/internal/pkg"
	"github.com/sirupsen/logrus"
	"os"
	"unsafe"
)

func Usage() {
	fmt.Println("mkfs.minix: Usage: mkfs.minix /dev/name")
	os.Exit(0)
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

//#define INODE_SIZE (sizeof(struct minix_inode))

const (
	MINIX_BLOCK_SIZE_BITS = int(10)
	MINIX_BLOCK_SIZE      = 1 << MINIX_BLOCK_SIZE_BITS
)

var (
	INODE_SIZE             = 0
	MINIX_INODES_PER_BLOCK = 0

	device_name = ""
	root_block  = [MINIX_BLOCK_SIZE]byte{}
)

func init() {
	INODE_SIZE = int(unsafe.Sizeof(pkg.Inode{}))
	MINIX_INODES_PER_BLOCK = (MINIX_BLOCK_SIZE) / int(unsafe.Sizeof(pkg.Inode{}))
}

func check_mount() {
	var err error
	mss, err := etc.GetMtab()
	if err != nil {
		logrus.Errorf("etc.GetMtab err:%v", err)
		panic(fmt.Sprintf("etc.GetMtab err:%v", err))
	}
	for _, ms := range mss {
		if ms.FileSystem == device_name {
			panic(fmt.Sprintf("%s is mounted; will not make a filesystem here!", device_name))
		}
	}

	return
}

func main() {
	if len(os.Args) != 2 {
		Usage()
	}
	device_name = os.Args[1]

	if INODE_SIZE != MINIX_BLOCK_SIZE {
		panic("bad inode size")
	}
	/*
		struct termios tmp;
		int count;
		int retcode = FSCK_EX_OK;

		setlocale(LC_ALL, ""); // 时区设置
		bindtextdomain(PACKAGE, LOCALEDIR);
		textdomain(PACKAGE);
		atexit(close_stdout);
		if (argc == 2 &&
		    (!strcmp(argv[1], "-V") || !strcmp(argv[1], "--version"))) {
			printf(UTIL_LINUX_VERSION);
			exit(FSCK_EX_OK);
		}

		if (INODE_SIZE * MINIX_INODES_PER_BLOCK != MINIX_BLOCK_SIZE)
			die(_("bad inode size"));
		if (INODE2_SIZE * MINIX2_INODES_PER_BLOCK != MINIX_BLOCK_SIZE)
			die(_("bad v2 inode size"));
			opterr = 0;
		while ((i = getopt(argc, argv, "ci:l:n:v123")) != -1)
			switch (i) {
			case 'c':
				check=1; break;
			case 'i':
				req_nr_inodes = strtoul_or_err(optarg,
						_("failed to parse number of inodes"));
				break;
			case 'l':
				listfile = optarg; break;
			case 'n':
				i = strtoul_or_err(optarg,
						_("failed to parse maximum length of filenames"));
				if (i == 14)
					magic = MINIX_SUPER_MAGIC;
				else if (i == 30)
					magic = MINIX_SUPER_MAGIC2;
				else
					usage();
				namelen = i;
				dirsize = i+2;
				break;
			case '1':
				fs_version = 1;
				break;
			case '2':
			case 'v': // kept for backwards compatiblitly
					fs_version = 2;
					break;
			case '3':
				fs_version = 3;
				namelen = 60;
				dirsize = 64;
				break;
			default:
				usage();
			}
		argc -= optind;
		argv += optind;
		if (argc > 0 && !device_name) {
			device_name = argv[0];
			argc--;
			argv++;
		}
		if (argc > 0)
			BLOCKS = strtoul_or_err(argv[0], _("failed to parse number of blocks"));

		if (!device_name) {
			usage();
		}
		check_mount();		// is it already mounted?

		tmp = root_block;
		if (fs_version == 3) {
			*(uint32_t *)tmp = 1;
			strcpy(tmp+4,".");
			tmp += dirsize;
			*(uint32_t *)tmp = 1;
			strcpy(tmp+4,"..");
			tmp += dirsize;
			*(uint32_t *)tmp = 2;
			strcpy(tmp+4, ".badblocks");
		} else {
			*(uint16_t *)tmp = 1;
			strcpy(tmp+2,".");
			tmp += dirsize;
			*(uint16_t *)tmp = 1;
			strcpy(tmp+2,"..");
			tmp += dirsize;
			*(uint16_t *)tmp = 2;
			strcpy(tmp+2, ".badblocks");
		}
		if (stat(device_name, &statbuf) < 0)
			err(MKFS_EX_ERROR, _("stat failed %s"), device_name);
		if (S_ISBLK(statbuf.st_mode))
			DEV = open(device_name,O_RDWR | O_EXCL);
		else
			DEV = open(device_name,O_RDWR);

		if (DEV<0)
			err(MKFS_EX_ERROR, _("cannot open %s"), device_name);

	*/

	root_block = InitRootBlock()

	check_mount()
}
