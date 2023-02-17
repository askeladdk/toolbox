// Package bloom provides a bloom filter implementation.
// A bloom filter is a space-efficient probabilistic data structure
// that quickly tests set membership given a certain probability
// of false positives but no false negatives.
// This implementation is optimized for modern L1 caches
// and is safe to use concurrently.
//
// A bloom filter is characterized by four interrelated parameters:
//
//   - k: The number of hash functions.
//   - m: The amount of memory in bits.
//   - n: The expected number of elements.
//   - p: The probability of false positives.
//
// In this implementation k is fixed at 8 and m is estimated given n and p.
//
// The bloom filter maintains a set of uint64 integers.
// The integers can be obtained by hashing a value using [hash.Hash64].
// or they can be added and tested directly.
// Note that integers should be mixed using [Uint64] if they are not already randomized,
// Otherwise the number of false positives will be unacceptably high.
package bloom

import (
	"math"
	"math/bits"
	"sync/atomic"
)

// References:
// Bloom filter calculator
// https://hur.st/bloomfilter
// Less Hashing, Same Performance: Building a Better Bloom Filter
// https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf
// Cache-, Hash- and Space-Efficient Bloom Filters
// https://algo2.iti.kit.edu/singler/publications/cacheefficientbloomfilters-wea2007.pdf
// Performance-Optimal Filtering: Bloom Overtakes Cuckoo at High Throughput
// https://www.vldb.org/pvldb/vol12/p502-lang.pdf

// Fixed parameters:
// k = 8
// S = 8 (64 bits)
// B = 512 (bits)
// z = 4

// Filter is a bloom filter.
type Filter []uint64

// New returns a new Filter having m bits of memory.
func New(m int) Filter {
	words := ((m + 511) / 512) * 8
	return make([]uint64, words)
}

// NewWithEstimate is shorthand for New(Estimate(n, p)).
func NewWithEstimate(n int, p float64) Filter {
	return New(Estimate(n, p))
}

// Add includes h in the filter.
// Adding h twice does not change the filter.
// The complexity is O(1).
func (f Filter) Add(x uint64) {
	h0, h1 := splithash(x)
	s0, s1, s2, s3 := sectors(h0, h1, uint32(len(f)))
	z0, z1, z2, z3 := bitmasks(h0, h1)
	f.atomicSetBits(s0, z0)
	f.atomicSetBits(s1, z1)
	f.atomicSetBits(s2, z2)
	f.atomicSetBits(s3, z3)
}

// Test reports whether h may be in the filter.
// Returns true if h probably exists
// and false if it definitely does not.
// The complexity is O(1).
func (f Filter) Test(x uint64) bool {
	h0, h1 := splithash(x)
	s0, s1, s2, s3 := sectors(h0, h1, uint32(len(f)))
	z0, z1, z2, z3 := bitmasks(h0, h1)
	m0 := atomic.LoadUint64(&f[s0])
	m1 := atomic.LoadUint64(&f[s1])
	m2 := atomic.LoadUint64(&f[s2])
	m3 := atomic.LoadUint64(&f[s3])
	return m0&z0 == z0 && m1&z1 == z1 && m2&z2 == z2 && m3&z3 == z3
}

// TestAndAdd is shorthand for Test(x) followed by Add(x)
// but is more efficient than calling them separately.
// The complexity is O(1).
func (f Filter) TestAndAdd(x uint64) bool {
	h0, h1 := splithash(x)
	s0, s1, s2, s3 := sectors(h0, h1, uint32(len(f)))
	z0, z1, z2, z3 := bitmasks(h0, h1)
	m0 := f.atomicSetBits(s0, z0)
	m1 := f.atomicSetBits(s1, z1)
	m2 := f.atomicSetBits(s2, z2)
	m3 := f.atomicSetBits(s3, z3)
	return m0&z0 == z0 && m1&z1 == z1 && m2&z2 == z2 && m3&z3 == z3
}

// Bits reports the number of bits in f.
func (f Filter) Bits() int {
	return len(f) * 64
}

// Empty reports whether f is empty.
// The complexity is O(n),
// but it is more efficient than testing for f.Len() == 0.
func (f Filter) Empty() bool {
	for i := range f {
		if atomic.LoadUint64(&f[i]) != 0 {
			return false
		}
	}

	return true
}

// Equal tests whether f and g are equal.
// The complexity is O(n).
func (f Filter) Equal(g Filter) bool {
	if len(f) != len(g) {
		return false
	}
	for i := range f {
		if atomic.LoadUint64(&f[i]) != atomic.LoadUint64(&g[i]) {
			return false
		}
	}
	return true
}

// Len estimates the number of elements in the filter.
// Returns [math.MaxInt] if all bits in f are set to 1,
// meaning that Test(x) is true for all x.
// The complexity is O(n).
func (f Filter) Len() int {
	// https://en.wikipedia.org/wiki/Bloom_filter#Approximating_the_number_of_items_in_a_Bloom_filter
	const k = 8
	const blockbits = 512
	var n float64
	for i := 0; i < len(f); i += 8 {
		var ones int
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+0]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+1]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+2]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+3]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+4]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+5]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+6]))
		ones += bits.OnesCount64(atomic.LoadUint64(&f[i+7]))
		if ones != 0 {
			n += math.Log1p(-float64(ones) / blockbits)
		}
	}

	if n == math.Inf(-1) {
		return math.MaxInt
	}

	return int(-(blockbits / k) * n)
}

// Reset clears the filter.
func (f Filter) Reset() {
	for i := range f {
		atomic.StoreUint64(&f[i], 0)
	}
}

// UnionWith sets f to the union of f and g.
// Panics if f.Bits() != g.Bits().
// The complexity is O(n).
func (f Filter) UnionWith(g Filter) {
	if len(f) != len(g) {
		panic("bloom: cannot union with filter of unequal size")
	}
	for i := range f {
		f.atomicSetBits(uint32(i), g[i])
	}
}

func (f Filter) atomicSetBits(s uint32, z uint64) (old uint64) {
	for {
		old = atomic.LoadUint64(&f[s])
		if old&z == z || atomic.CompareAndSwapUint64(&f[s], old, old|z) {
			return
		}
	}
}

// Estimate calculates the number of bits m
// based on the expected number of elements n and false positive rate p.
func Estimate(n int, p float64) (m int) {
	// https://hur.st/bloomfilter
	// n = ceil(m / (-k / log(1 - exp(log(p) / k))))
	// p = pow(1 - exp(-k / (m / n)), k)
	// m = ceil((n * log(p)) / log(1 / pow(2, log(2))));
	// k = round((m / n) * log(2));
	const c = -0.4804530139182015 // log(1/pow(2,log(2)))
	m = int(math.Ceil(float64(n) * math.Log(p) / c))
	m = (m + 511) &^ 511
	return
}

// Uint64 hashes x using the moremur mixer function.
// Use this to pseudo-randomize sequential integers
// so that they can be added to the bloom filter.
// Note that Uint64(0) = 0.
func Uint64(x uint64) uint64 {
	// https://mostlymangling.blogspot.com/2019/12/stronger-better-morer-moremur-better.html
	x ^= x >> 27
	x *= 0x3C79AC492BA7B653
	x ^= x >> 33
	x *= 0x1C69B3F74AC4AE35
	x ^= x >> 27
	return x
}

// Int is shorthand for Uint64(uint64(x)).
func Int(x int) uint64 {
	return Uint64(uint64(x))
}

func sectors(h0, h1, n uint32) (s0, s1, s2, s3 uint32) {
	bl := (h0 % n) &^ 7
	s0 = bl + 0 + (1*h1)>>31
	s1 = bl + 2 + (2*h1)>>31
	s2 = bl + 4 + (3*h1)>>31
	s3 = bl + 6 + (4*h1)>>31
	return
}

func bitmasks(h0, h1 uint32) (z0, z1, z2, z3 uint64) {
	z0 |= 1 << (((h0 + 1*h1) >> 16) & 63)
	z0 |= 1 << (((h0 + 2*h1) >> 16) & 63)
	z1 |= 1 << (((h0 + 3*h1) >> 16) & 63)
	z1 |= 1 << (((h0 + 4*h1) >> 16) & 63)
	z2 |= 1 << (((h0 + 5*h1) >> 16) & 63)
	z2 |= 1 << (((h0 + 6*h1) >> 16) & 63)
	z3 |= 1 << (((h0 + 7*h1) >> 16) & 63)
	z3 |= 1 << (((h0 + 8*h1) >> 16) & 63)
	return
}

func splithash(h uint64) (h0, h1 uint32) {
	return uint32(h), uint32(h >> 32)
}
