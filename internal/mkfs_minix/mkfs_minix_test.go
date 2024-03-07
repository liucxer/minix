package mkfs_minix_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/minix/internal/mkfs_minix"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitInodeBitMap(t *testing.T) {
	superBlock := mkfs_minix.InitSuperBlock()
	res := mkfs_minix.InitInodeBitMap(superBlock)
	spew.Dump(res.Count())
	spew.Dump(res.Bytes())
	bts, err := res.MarshalBinary()
	require.NoError(t, err)
	spew.Dump(bts)
}
