package sparsebits

import (
	"math/rand"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestSparseSetUnset(t *testing.T) {
	n := 32000
	b := New(n)

	require.True(t, b.Len()&(b.Len()-1) == 0)

	require.True(t, !b.TestBit(0))
	b.SetBit(0, false)

	rnd := rand.New(rand.NewSource(0))
	xs := rnd.Perm(b.Len())

	for _, x := range xs {
		b.SetBit(x, true)
		require.True(t, b.TestBit(x))
	}

	for _, x := range xs {
		x = xs[x]
		b.FlipBit(x)
		require.True(t, !b.TestBit(x))
	}

	require.Equal(t, 0, b.OnesCount())

	require.True(t, !b.TestBit(b.Len()))
}

func TestSparseCount(t *testing.T) {
	b := New(1000)

	rnd := rand.New(rand.NewSource(0))
	xs := rnd.Perm(b.Len())

	for _, x := range xs {
		b.SetBit(x, true)
	}

	for i := 0; i < len(xs); i += 2 {
		b.SetBit(i, false)
	}

	require.Equal(t, b.Len()/2, b.OnesCount())

	b.Reset()
	require.Equal(t, 0, b.OnesCount())

	b.SetBit(65, true)
	require.Equal(t, 1, b.OnesCount())
}
