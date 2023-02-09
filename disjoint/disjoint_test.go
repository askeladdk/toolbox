package disjoint

import (
	"math/rand"
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestSet(t *testing.T) {
	set := New(make([]int, 8))
	set.Union(0, 0)
	set.Union(0, 1)
	set.Union(0, 2)
	set.Union(3, 3)
	set.Union(3, 4)
	set.Union(5, 5)
	set.Union(5, 6)
	set.Union(7, 7)

	require.Equal(t, 8, set.Len(), "Len")
	require.Equal(t, 4, set.CountGroups(), "CountGroups")
	require.True(t, set.Same(0, 1), "Same")
	require.True(t, !set.Same(0, 7), "not Same")

	roots := []int{1, 1, 1, 4, 4, 6, 6, 7}
	sizes := []int{3, 3, 3, 2, 2, 2, 2, 1}

	for i := range set {
		require.Equal(t, roots[i], set.Find(i), "Find", i)
		require.Equal(t, sizes[i], set.Size(i), "Size", i)
	}
}

func TestZeroPanic(t *testing.T) {
	invalid := make(Set, 10)

	var panicked bool

	func() {
		defer func() {
			if v := recover(); v != nil {
				panicked = true
			}
		}()
		invalid.Find(4)
	}()

	require.True(t, panicked, "Panicked")
}

func TestStressTest(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))

	set := New(make([]int, 10000))

	is := rnd.Perm(set.Len())
	for i := 1; i < len(is); i++ {
		set.Union(is[i-1], is[i])
	}

	require.Equal(t, 1, set.CountGroups(), "CountGroups")
	require.Equal(t, 10000, set.Size(0), "Size")
}
