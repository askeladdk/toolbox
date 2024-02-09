//go:build !amd64

package murmurhash3

import "encoding/binary"

func getblock(p []byte) (uint64, uint64) {
	_ = p[15]
	k0 := binary.LittleEndian.Uint64(p[0:8])
	k1 := binary.LittleEndian.Uint64(p[8:16])
	return k0, k1
}
