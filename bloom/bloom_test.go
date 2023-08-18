package bloom

import (
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestAddTest(t *testing.T) {
	f := New(1000000)
	require.True(t, f.Bits()%512 == 0)
	f.Add(Uint64(1))
	require.True(t, f.Test(Uint64(1)))
	require.True(t, !f.Empty())
	f.Reset()
	require.True(t, f.Empty())
}

func TestEqual(t *testing.T) {
	f := New(100000)
	g := New(100000)
	h := New(200000)
	require.True(t, f.Equal(g))

	for i := 1; i <= 1000; i++ {
		f.Add(Int(i))
		g.Add(Int(i))
		h.Add(Int(i))
	}

	require.True(t, f.Equal(g))
	g.Add(Uint64(1337))
	require.True(t, !f.Equal(g))
	require.True(t, !f.Equal(h))

	require.True(t, f.Len()-1000 < 10)
	require.True(t, h.Len()-1000 < 10)
}

func TestFalsePositiveRate(t *testing.T) {
	p := 0.001
	n := 1000000
	f := NewWithEstimate(n, p)

	var collisions int
	for i := 1; i <= n; i++ {
		if f.TestAndAdd(Int(i)) {
			collisions++
		}
	}

	require.True(t, float64(collisions)/float64(n) <= p)
	require.True(t, f.Len()-n < 1000)
}

func TestUnion(t *testing.T) {
	f := New(1000000)
	g := New(1000000)

	for i := 1; i <= 10000; i++ {
		f.Add(Int(i))
		g.Add(Int(10000 + i))
	}

	n := f.Len() + g.Len()
	f.UnionWith(g)

	d := n - f.Len()
	if d < 0 {
		d = -d
	}

	require.True(t, d < 10)
}

func TestLenMax(t *testing.T) {
	f := New(100000)
	for i := 0; i < len(f); i++ {
		f[i] = math.MaxUint64
	}
	require.Equal(t, math.MaxInt, f.Len())
}

func TestUnionPanic(t *testing.T) {
	f := New(1000000)
	g := New(2000000)

	var panicked bool

	func() {
		defer func() {
			panicked = recover() != nil
		}()
		f.UnionWith(g)
	}()

	require.True(t, panicked)
}

func TestParallelism(t *testing.T) {
	n := 1000000

	g := NewWithEstimate(n, 0.001)
	for i := 1; i <= n; i++ {
		g.Add(Int(i))
	}

	for _, c := range []int{2, 4, 8, 16, 32, 64} {
		t.Run(strconv.Itoa(c), func(t *testing.T) {
			f := NewWithEstimate(n, 0.001)
			d := n / c
			var wg sync.WaitGroup

			for i := 1; i <= n; i += d {
				wg.Add(1)
				go func(i, j int) {
					for ; i < j; i++ {
						f.Add(Int(i))
					}
					wg.Done()
				}(i, i+d)
			}

			wg.Wait()
			require.True(t, f.Equal(g))
		})
	}
}

func BenchmarkBloomFilter(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000, 1000000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			f := NewWithEstimate(n, 0.01)
			var fp int
			for i := 1; i <= n; i++ {
				if f.TestAndAdd(Int(i)) {
					fp++
				}
			}
			b.ReportMetric(float64(fp)/float64(n), "fpr")
		})
	}
}

func BenchmarkParallelism(b *testing.B) {
	f := New(Estimate(b.N, 0.01))
	b.ResetTimer()

	c := b.N / runtime.GOMAXPROCS(0)
	var i uint64

	b.RunParallel(func(p *testing.PB) {
		j := atomic.AddUint64(&i, 1) * uint64(c)
		for ; p.Next(); j++ {
			f.Add(Uint64(j))
		}
	})
}
