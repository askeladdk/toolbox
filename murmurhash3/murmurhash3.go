// Package murmurhash3 provides the MurmurHash3-x64-128 hash function.
// Murmurhash3 is a fast non-cryptographic hash function suitable for general hash-based lookup.
// It was created by Austin Appleby in 2008. See: https://github.com/aappleby/smhasher
package murmurhash3

import (
	"encoding/binary"
	"errors"
	"hash"
	"math/bits"
)

// Size in bytes of a MurmurHash3 checksum.
const Size = 16

const (
	c0 = 0x87c37b91114253d5
	c1 = 0x4cf5ad432745937f
)

func round(h0, h1, k0, k1 uint64) (uint64, uint64) {
	// h0
	k0 *= c0
	k0 = bits.RotateLeft64(k0, 31)
	k0 *= c1
	h0 ^= k0
	h0 = bits.RotateLeft64(h0, 27)
	h0 += h1
	h0 = h0*5 + 0x52dce729
	// h1
	k1 *= c1
	k1 = bits.RotateLeft64(k1, 33)
	k1 *= c0
	h1 ^= k1
	h1 = bits.RotateLeft64(h1, 31)
	h1 += h0
	h1 = h1*5 + 0x38495ab5
	return h0, h1
}

func mix(h0, h1 uint64, p []byte) (uint64, uint64, []byte) {
	for ; len(p) >= 16; p = p[16:] {
		k0, k1 := getblock(p)
		h0, h1 = round(h0, h1, k0, k1)
	}
	return h0, h1, p
}

func fmix64(h uint64) uint64 {
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

func finalize(h0, h1, n, head uint64, tail [16]byte) (uint64, uint64) {
	var k0, k1 uint64

	switch head & 15 { // head is always <= 15
	case 15:
		k1 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k1 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k1 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k1 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k1 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k1 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k1 ^= uint64(tail[8])
		k1 *= c1
		k1 = bits.RotateLeft64(k1, 33)
		k1 *= c0
		h1 ^= k1
		fallthrough
	case 8:
		k0 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k0 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k0 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k0 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k0 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k0 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k0 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k0 ^= uint64(tail[0])
		k0 *= c0
		k0 = bits.RotateLeft64(k0, 31)
		k0 *= c1
		h0 ^= k0
	}

	h0 ^= n
	h1 ^= n
	h0 += h1
	h1 += h0
	h0 = fmix64(h0)
	h1 = fmix64(h1)
	h0 += h1
	h1 += h0

	return h0, h1
}

type digest struct {
	s0   uint64
	s1   uint64
	h0   uint64
	h1   uint64
	n    uint64
	head uint64
	tail [16]byte
}

type digest128 struct {
	digest
}

type digest64 struct {
	digest
}

func (d *digest) BlockSize() int {
	return 1
}

func (d *digest128) Size() int {
	return 16
}

func (d *digest64) Size() int {
	return 8
}

func (d *digest) Reset() {
	d.h0 = d.s0
	d.h1 = d.s1
	d.n = 0
	d.head = 0
	clear(d.tail[:])
}

func (d *digest) Write(p []byte) (int, error) {
	d.n += uint64(len(p))

	h0, h1 := d.h0, d.h1

	if d.head > 0 {
		r := 16 - d.head
		if uint64(len(p)) < r {
			copy(d.tail[d.head:], p)
			d.head += uint64(len(p))
			return len(p), nil
		}

		copy(d.tail[d.head:], p[:r])
		p = p[r:]
		k0, k1 := getblock(d.tail[:])
		h0, h1 = round(h0, h1, k0, k1)
		d.head = 0
	}

	var tail []byte
	d.h0, d.h1, tail = mix(h0, h1, p)
	copy(d.tail[0:], tail)
	d.head = uint64(len(tail))
	return len(p), nil
}

func (d *digest128) Sum(b []byte) []byte {
	h0, h1 := finalize(d.h0, d.h1, d.n, d.head, d.tail)
	b = binary.BigEndian.AppendUint64(b, h0)
	b = binary.BigEndian.AppendUint64(b, h1)
	return b
}

func (d *digest64) Sum(b []byte) []byte {
	h0, _ := finalize(d.h0, d.h1, d.n, d.head, d.tail)
	b = binary.BigEndian.AppendUint64(b, h0)
	return b
}

func (d *digest64) Sum64() uint64 {
	h0, _ := finalize(d.h0, d.h1, d.n, d.head, d.tail)
	return h0
}

func (d *digest) MarshalBinary() ([]byte, error) {
	b := make([]byte, 64)
	binary.BigEndian.PutUint64(b[0:], d.s0)
	binary.BigEndian.PutUint64(b[8:], d.s1)
	binary.BigEndian.PutUint64(b[16:], d.h0)
	binary.BigEndian.PutUint64(b[24:], d.h1)
	binary.BigEndian.PutUint64(b[32:], d.n)
	binary.BigEndian.PutUint64(b[40:], d.head)
	copy(b[48:], d.tail[0:])
	return b, nil
}

func (d *digest) UnmarshalBinary(b []byte) error {
	if len(b) != 64 {
		return errors.New("murmurhash3: invalid hash state size")
	}
	d.s0 = binary.BigEndian.Uint64(b[0:])
	d.s1 = binary.BigEndian.Uint64(b[8:])
	d.h0 = binary.BigEndian.Uint64(b[16:])
	d.h1 = binary.BigEndian.Uint64(b[24:])
	d.n = binary.BigEndian.Uint64(b[32:])
	d.head = binary.BigEndian.Uint64(b[40:])
	copy(d.tail[0:], b[48:])
	return nil
}

// New returns a new [hash.Hash] initialized with a zero seed.
func New() hash.Hash {
	return &digest128{}
}

// NewWithSeed returns a new [hash.Hash] initialized with the given seed.
func NewWithSeed(s0, s1 uint64) hash.Hash {
	var d digest128
	d.s0 = s0
	d.s1 = s1
	d.h0 = s0
	d.h1 = s1
	return &d
}

// New64 returns a new [hash.Hash64] initialized with a zero seed.
func New64() hash.Hash64 {
	return &digest64{}
}

// New64WithSeed returns a new [hash.Hash64] initialized with the given seed.
func New64WithSeed(s0, s1 uint64) hash.Hash64 {
	var d digest64
	d.s0 = s0
	d.s1 = s1
	d.h0 = s0
	d.h1 = s1
	return &d
}

// Sum calculates the 128-bit hash of p.
func Sum(p []byte) [Size]byte {
	var b [Size]byte
	var d digest128
	_, _ = d.Write(p)
	d.Sum(b[:0])
	return b
}

// Sum64 calculates the 64-bit hash of p.
func Sum64(p []byte) uint64 {
	var d digest64
	_, _ = d.Write(p)
	return d.Sum64()
}
