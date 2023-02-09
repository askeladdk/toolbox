// Package queue provides a generic queue implementation.
// A queue is a first-in-first-out data structure
// in which elements are removed in the order they were added.
package queue

// Queue implements a generic queue.
// A queue is a first-in-first-out data structure
// in which elements are removed in the order they were added.
//
// The queue acts as a ring buffer as long as there is capacity for more elements,
// but resizes when it is at capacity and a new element is pushed.
//
// The zero value for Queue is an empty queue ready to use.
type Queue[T any] struct {
	elem []T
	head int
	tail int
}

// New returns a new empty queue with an initial capacity
// provided by buf.
func New[T any](buf []T) *Queue[T] {
	return &Queue[T]{
		elem: buf,
	}
}

// Cap reports the total capacity of q.
func (q *Queue[T]) Cap() int {
	return len(q.elem)
}

// Empty reports whether q is empty.
func (q *Queue[T]) Empty() bool {
	return q.head == q.tail
}

// Len reports the number of elements in q.
func (q *Queue[T]) Len() int {
	if len(q.elem) == 0 {
		return 0
	}
	return q.wraparound(q.head - q.tail)
}

// Peek returns the front element of q without removing it.
// Panics if q is empty.
func (q *Queue[T]) Peek() T {
	if q.Empty() {
		panic("queue: underflow")
	}

	return q.elem[q.index(q.tail)]
}

// Peek removes and returns the front element of q.
// Panics if q is empty.
func (q *Queue[T]) Pop() T {
	var zeroval T
	x := q.Peek()
	q.elem[q.index(q.tail)] = zeroval
	q.tail = q.wraparound(q.tail + 1)
	return x
}

// Push appends x to the front of q.
func (q *Queue[T]) Push(x T) {
	q.Grow(1)
	q.elem[q.index(q.head)] = x
	q.head = q.wraparound(q.head + 1)
}

// Grow ensures that q has capacity for at least n elements.
func (q *Queue[T]) Grow(n int) {
	if n < 0 {
		panic("queue: cannot negative grow")
	}

	size := len(q.elem)

	if size == 0 {
		q.elem = make([]T, n)
		return
	}

	if size-q.wraparound(q.head-q.tail) >= n {
		return
	}

	elem := append(q.elem, make([]T, n)...)
	elem = elem[:cap(elem)]

	// if the index of head <= tail,
	// we have swap the elements so that head
	// comes after tail and new elements can be pushed
	// without overwriting the oldest.
	//
	// Before:
	// hhhhhhhhttttttttttttt
	// ^      ^^
	// 0  head  tail
	//
	// After:
	// ttttttttttttthhhhhhhh--------------------
	// ^                   ^                   ^
	// 0 tail             head                cap

	hi := q.index(q.head)
	lo := q.index(q.tail)

	if hi <= lo {
		// copy the shortest run to tmp
		if hi < size-lo {
			tmp := make([]T, hi)
			copy(tmp, elem[:hi])
			copy(elem, elem[lo:size])
			copy(elem[size-lo:], tmp)
		} else {
			tmp := make([]T, size-lo)
			copy(tmp, elem[lo:size])
			copy(elem[size-lo:], elem[:hi])
			copy(elem, tmp)
		}

		q.head = size
		q.tail = 0
	}

	q.elem = elem
}

func (q *Queue[T]) index(i int) int {
	return i % len(q.elem)
}

func (q *Queue[T]) wraparound(i int) int {
	return i % (2 * len(q.elem))
}
