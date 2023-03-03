// Package densebits provides operations on dense bit sets.
package densebits

import (
	"math/bits"
	"strconv"
	"strings"
)

// Set represents a dense bit set.
// It is implemented as a slice of unsigned integers.
// Use it to bitwise manipulate a range of bits.
type Set []uint64

// New returns a new dense bit set for bits in range [0, n)
// where n is rounded up to the nearest multiple of 64.
func New(n int) Set {
	return make(Set, (n+63)/64)
}

// Len returns the number of bits in s.
func (s Set) Len() int {
	return len(s) * 64
}

// Set sets or clears the i-th bit.
// The complexity is O(1).
func (s Set) Set(i int, to bool) {
	w, b := wordbit(i)
	if to {
		s[w] |= b
	} else {
		s[w] &= ^b
	}
}

// Get reports whether the i-th bit is set to one.
// The complexity is O(1).
func (s Set) Get(i int) bool {
	w, b := wordbit(i)
	return s[w]&b != 0
}

// Flip sets the i-th bit to one if it zero or to zero it if is one.
// The complexity is O(1).
func (s Set) Flip(i int) {
	w, b := wordbit(i)
	s[w] ^= b
}

// OnesCount reports the number of one bits (population count) in s.
// The complexity is O(n).
func (s Set) OnesCount() int {
	count := 0
	for i := range s {
		count += bits.OnesCount64(s[i])
	}
	return count
}

// Equal reports whether s and p are equal.
// The complexity is O(n).
func (s Set) Equal(p Set) bool {
	if len(s) != len(p) {
		return false
	}
	for i := range s {
		if s[i] != p[i] {
			return false
		}
	}
	return true
}

// And stores the result of p AND q in s,
// which will be resized to min(|p|, |q|) bits if needed.
// The complexity is O(n).
func (s *Set) And(p, q Set) {
	m := min(len(p), len(q))
	s.sizeto(m)
	for i := 0; i < m; i++ {
		(*s)[i] = p[i] & q[i]
	}
}

// Or stores the result of p OR q in s,
// which will be resized to min(|p|, |q|) bits if needed.
// The complexity is O(n).
func (s *Set) Or(p, q Set) {
	m := min(len(p), len(q))
	s.sizeto(m)
	for i := 0; i < m; i++ {
		(*s)[i] = p[i] | q[i]
	}
}

// AndNot stores the result of p AND NOT q in s,
// which will be resized to min(|p|, |q|) bits if needed.
// The complexity is O(n).
func (s *Set) AndNot(p, q Set) {
	m := min(len(p), len(q))
	s.sizeto(m)
	for i := 0; i < m; i++ {
		(*s)[i] = p[i] &^ q[i]
	}
}

// Xor stores the result of p XOR q in s,
// which will be resized to min(|p|, |q|) bits if needed.
// The complexity is O(n).
func (s *Set) Xor(p, q Set) {
	m := min(len(p), len(q))
	s.sizeto(m)
	for i := 0; i < m; i++ {
		(*s)[i] = p[i] ^ q[i]
	}
}

// Not stores the result of NOT p in s,
// which will be resized to |p| bits if needed.
// The complexity is O(n).
func (s *Set) Not(p Set) {
	s.sizeto(len(p))
	for i := range *s {
		(*s)[i] = ^p[i]
	}
}

// ShiftLeft stores the result of p << n in s,
// where n is a number from 0 to 64,
// and s will be resized to |p| bits if needed.
// The complexity is O(n).
func (s *Set) ShiftLeft(p Set, n int) (remainder uint64) {
	m := min(len(*s), len(p))
	s.sizeto(m)
	remainder = p[0] >> uint64(64-n)
	(*s)[0] = p[0] << n
	for i := 1; i < m; i++ {
		pi := p[i]
		(*s)[i-1] |= (pi >> uint64(64-n))
		(*s)[i] = pi << n
	}
	return
}

// ShiftRight stores the result of p >> n in s,
// where n is a number from 0 to 64,
// and s will be resized to |p| bits if needed.
// The complexity is O(n).
func (s *Set) ShiftRight(p Set, n int) (remainder uint64) {
	m := min(len(*s), len(p))
	s.sizeto(m)
	mask := (uint64(1) << n) - 1
	for i := 0; i < m; i++ {
		pi := p[i]
		x := (pi >> n) | remainder
		remainder = (pi & mask) << (64 - n)
		(*s)[i] = x
	}
	return
}

// RotateLeft stores the result of p left rotated by n bits in s,
// where n is a number from 0 to 64,
// and s will be resized to |p| bits if needed.
// The complexity is O(n).
func (s *Set) RotateLeft(p Set, n int) {
	remainder := s.ShiftLeft(p, n)
	(*s)[len(*s)-1] |= remainder
}

// RotateRight stores the result of p right rotated by n bits in s,
// where n is a number from 0 to 64,
// and s will be resized to |p| bits if needed.
// The complexity is O(n).
func (s *Set) RotateRight(p Set, n int) {
	remainder := s.ShiftRight(p, n)
	(*s)[0] |= remainder
}

// Fill sets all words to x.
// The complexity is O(n).
func (s Set) Fill(x uint64) {
	for i := 0; i < len(s); i++ {
		s[i] = x
	}
}

// Reset shorthand for Fill(0).
// The complexity is O(n).
func (s Set) Reset() {
	s.Fill(0)
}

// Slice returns a slice of s in the range of [lo, hi) bits,
// where lo and hi are rounded to a multiple of 64 bits.
// The slice obeys all the ordinary slice rules.
func (s Set) Slice(lo, hi int) Set {
	return s[lo/64 : (hi+63)/64]
}

// String implements fmt.Stringer.
func (s Set) String() string {
	var b [16]byte
	var sb strings.Builder
	sb.Grow(2 + 16*len(s) + len(s) - 1)
	sb.WriteByte('[')
	if len(s) > 0 {
		sb.Write(strconv.AppendUint(b[:], s[0], 16))
		for _, x := range s[1:] {
			sb.WriteByte(' ')
			sb.Write(strconv.AppendUint(b[:], x, 16))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

func (s *Set) sizeto(n int) {
	if len(*s) < n {
		*s = append(*s, make(Set, n-len(*s))...)
		return
	}
	*s = (*s)[:n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func wordbit(i int) (w int, b uint64) {
	return i / 64, 1 << (i % 64)
}
