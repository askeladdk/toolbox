//go:build amd64

package murmurhash3

import "unsafe"

func getblock(p []byte) (uint64, uint64) {
	_ = p[15]
	k0 := *(*uint64)(unsafe.Pointer(&p[0]))
	k1 := *(*uint64)(unsafe.Pointer(&p[8]))
	return k0, k1
}
