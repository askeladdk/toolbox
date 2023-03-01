package sparse

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestMap(t *testing.T) {
	m := NewMap[uint, int](0)
	m.Grow(1)

	require.True(t, m.Cap() >= 1)

	rnd := rand.New(rand.NewSource(0))
	perms := rnd.Perm(100)

	for i := 0; i < 50; i++ {
		m.Set(uint(perms[i]), -perms[i])
	}

	for i := 0; i < 50; i++ {
		x := uint(perms[i])
		require.True(t, m.Has(x))
		j, ok := m.Index(x)
		require.True(t, ok)
		require.Equal(t, -perms[i], m.Values()[j])
		v, ok := m.Get(x)
		require.True(t, ok)
		require.Equal(t, -perms[i], v)
	}

	require.Equal(t, 50, len(m.Keys()))
	require.Equal(t, 50, len(m.Values()))

	sort.Sort(m)
	require.True(t, sort.IsSorted(m))

	for _, x := range perms {
		m.Del(uint(x))
		require.True(t, !m.Has(uint(x)))
	}

	require.Equal(t, 0, m.Len())
	m.Set(1, 1)
	m.Reset()
	require.Equal(t, 0, m.Len())
	require.True(t, m.Empty())
}
