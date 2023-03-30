package queue

import (
	"testing"

	"github.com/askeladdk/toolbox/internal/require"
)

func TestPushPop(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	q := New[int](10)

	require.Equal(t, 0, q.Len())

	for _, v := range values {
		q.Push(v)
	}

	require.Equal(t, len(values), q.Len())

	for i := 0; !q.Empty(); i++ {
		require.Equal(t, values[i], q.Pop())
	}

	require.True(t, q.Empty())
	require.Equal(t, 0, q.Len())
	require.True(t, q.Cap() >= len(values))
}

func TestPushPopGrow(t *testing.T) {
	var q Queue[int]

	for i := 1; i <= 10000; i++ {
		q.Push(i)
	}

	for i := 1; !q.Empty(); i++ {
		require.Equal(t, i, q.Pop())
	}
}

func TestPushPopWraparound(t *testing.T) {
	q := Queue[int]{
		elem: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		head: 10,
	}

	require.Equal(t, 10, q.Len())

	for i := 1; i <= 5; i++ {
		require.Equal(t, i, q.Pop())
	}

	require.Equal(t, 5, q.Len())

	for i := 11; i <= 15; i++ {
		q.Push(i)
	}

	require.Equal(t, 10, q.Len())
	require.Equal(t, 10, q.Cap())
	require.Equal(t, []int{11, 12, 13, 14, 15, 6, 7, 8, 9, 10}, q.elem)
	require.Equal(t, 10+5, q.head)
	require.Equal(t, 5, q.tail)
}

func TestGrowReorderTmpHead(t *testing.T) {
	q := Queue[int]{
		elem: []int{7, 8, 9, 10, 1, 2, 3, 4, 5, 6},
		head: 10 + 4,
		tail: 4,
	}

	require.Equal(t, 10, q.Len())

	q.Push(11)

	require.Equal(t, 11, q.Len())
	require.True(t, q.Cap() > 10)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, q.elem[:q.Len()])
	require.Equal(t, 0, q.tail)
	require.Equal(t, 11, q.head)

	for i := 1; !q.Empty(); i++ {
		require.Equal(t, i, q.Pop())
	}
}

func TestGrowReorderTmpTail(t *testing.T) {
	q := Queue[int]{
		elem: []int{5, 6, 7, 8, 9, 10, 1, 2, 3, 4},
		head: 10 + 6,
		tail: 6,
	}

	require.Equal(t, 10, q.Len())

	q.Push(11)

	require.Equal(t, 11, q.Len())
	require.True(t, q.Cap() > 10)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, q.elem[:q.Len()])
	require.Equal(t, 0, q.tail)
	require.Equal(t, 11, q.head)

	for i := 1; !q.Empty(); i++ {
		require.Equal(t, i, q.Pop())
	}
}

func TestGrowNegPanic(t *testing.T) {
	var q Queue[int]
	var panicked bool

	func() {
		defer func() {
			if v := recover(); v != nil {
				panicked = true
			}
		}()

		q.Grow(-1)
	}()

	require.True(t, panicked)
}

func TestPeekEmptyPanic(t *testing.T) {
	var q Queue[int]
	var panicked bool

	func() {
		defer func() {
			if v := recover(); v != nil {
				panicked = true
			}
		}()

		q.Peek()
	}()

	require.True(t, panicked)
}

func TestPushPopPush(t *testing.T) {
	var q Queue[int]
	require.Equal(t, 0, q.Len())

	for i := 0; i < 1000; i++ {
		q.Push(i)
		q.Pop()
		q.Push(i)
	}

	require.Equal(t, 1000, q.Len())

	for i := 500; i < 1000; i++ {
		require.Equal(t, i, q.Pop())
		require.Equal(t, i, q.Pop())
	}

	q.Reset()
	require.Equal(t, 0, q.Len())
}
