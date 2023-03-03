// Package sparsebits provides operations on sparse bit sets.
package sparsebits

import (
	"math/bits"
)

// Based on:
// https://bugfix-66.com/7256e0772dc3b02d72abf15b171731c933fd44d67de074d679f1e4cb7bb20f79

// Set represents a sparse bit set.
// It is implemented as a binary tree and is best suited when few bits are expected
// to be set to one and set membership is the only needed operation.
type Set struct {
	dat []uint64
	mid uint64
}

// New returns a new sparse bitset for bits in range [0, n)
// where n is rounded up to the nearest power of two.
func New(n int) *Set {
	have := uint64(128)
	for have < uint64(n) {
		have <<= 1
	}

	return &Set{
		dat: []uint64{0, 0},
		mid: have >> 1,
	}
}

// Len reports the number of bits in s.
func (s *Set) Len() int {
	return int(s.mid << 1)
}

// Reset clears the bitset and reuses the backing slice.
// The complexity is O(1).
func (s *Set) Reset() {
	s.dat = append(s.dat[:0], 0, 0)
}

// Set sets or clears the i-th bit.
// The complexity is O(log(n)).
func (s *Set) Set(i int, to bool) {
	if to {
		s.twiddle(i, 1)
	} else {
		s.twiddle(i, 0)
	}
}

// Flip sets the i-th bit to one if it zero or to zero it if is one.
// The complexity is O(log(n)).
func (s *Set) Flip(i int) {
	s.twiddle(i, 2)
}

// OnesCount reports the number of one bits (population count) in s.
// The complexity is O(n).
func (s *Set) OnesCount() int {
	var count int
	q := make([]uint64, 0, bits.TrailingZeros64(s.mid)<<1-2)
	q = append(q, 0, s.mid)

	for n := 2; n != 0; n = len(q) {
		mid := q[n-1]
		at := q[n-2]
		q = q[:n-2]

		at1 := s.dat[at|1]
		at = s.dat[at]

		if mid == 64 {
			count += bits.OnesCount64(at)
			count += bits.OnesCount64(at1)
		} else {
			mid >>= 1
			if at != 0 {
				q = append(q, at, mid)
			}
			if at1 != 0 {
				q = append(q, at1, mid)
			}
		}
	}

	return count
}

// Get reports whether the i-th bit is set to one.
// The complexity is O(log(n)).
func (s *Set) Get(i int) bool {
	idx := uint64(i)
	at := uint64(0)
	mid := s.mid
	for ; ; mid >>= 1 {
		if idx >= mid {
			idx -= mid
			at++
		}

		at = s.dat[at]

		if mid == 64 {
			return at&(1<<idx) != 0
		} else if at == 0 {
			return false
		}
	}
}

func (s *Set) twiddle(i, op int) {
	idx := uint64(i)
	at := uint64(0)
	mid := s.mid
	for ; ; mid >>= 1 {
		if idx >= mid {
			idx -= mid
			at++
		}

		if mid == 64 {
			bit := uint64(1) << idx
			switch op {
			case 0: // clear
				s.dat[at] &^= bit
			case 1: // set
				s.dat[at] |= bit
			case 2: // flip
				s.dat[at] ^= bit
			}
			return
		}

		down := s.dat[at]
		if down == 0 {
			if op == 0 { // clear
				return
			}
			down = uint64(len(s.dat))
			s.dat[at] = down
			s.dat = append(s.dat, 0, 0)
		}
		at = down
	}
}
