// Package queue provides a generic queue implementation.
// A queue is a first-in-first-out data structure
// in which elements are removed in the order they were added.
package queue

// Queue implements a generic queue.
//
// The zero value for Queue is an empty queue ready to use.
type Queue[E any] struct {
	elem []E
	head int
	tail int
}

// New returns a new empty queue with an initial capacity.
func New[E any](cap int) *Queue[E] {
	return &Queue[E]{
		elem: make([]E, cap),
	}
}

// Cap reports the total capacity of q.
func (q *Queue[E]) Cap() int {
	return len(q.elem)
}

// Empty reports whether q is empty.
// The complexity is O(1).
func (q *Queue[E]) Empty() bool {
	return q.head == q.tail
}

// Reset clears the queue without deallocation the underlying memory.
func (q *Queue[E]) Reset() {
	q.head = 0
	q.tail = 0
	// clear any possible pointers to prevent memory leaks
	clear(q.elem)
}

// Len reports the number of elements in q.
// The complexity is O(1).
func (q *Queue[E]) Len() int {
	if len(q.elem) == 0 {
		return 0
	}
	return q.wraparound(q.head - q.tail)
}

// Peek returns the front element of q without removing it.
// Panics if q is empty.
// The complexity is O(1).
func (q *Queue[E]) Peek() E {
	if q.Empty() {
		panic("queue: underflow")
	}

	return q.elem[q.index(q.tail)]
}

// Peek removes and returns the front element of q.
// Panics if q is empty.
// The complexity is O(1).
func (q *Queue[E]) Pop() E {
	var zeroval E
	x := q.Peek()
	q.elem[q.index(q.tail)] = zeroval
	q.tail = q.wraparound(q.tail + 1)
	return x
}

// Push appends x to the front of q
// and automatically allocates more capacity when needed.
// The complexity is O(1) amortized.
func (q *Queue[E]) Push(x E) {
	q.Grow(1)
	q.elem[q.index(q.head)] = x
	q.head = q.wraparound(q.head + 1)
}

// Grow ensures that q has capacity for at least n elements.
// The complexity is O(n) if resizing is needed.
func (q *Queue[E]) Grow(n int) {
	if n < 0 {
		panic("queue: cannot negative grow")
	}

	size := len(q.elem)

	if size == 0 {
		q.elem = make([]E, n)
		return
	}

	if size-q.wraparound(q.head-q.tail) >= n {
		return
	}

	elem := q.elem
	elem = append(elem, make([]E, n)...)
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
			tmp := make([]E, hi)
			copy(tmp, elem[:hi])
			copy(elem, elem[lo:size])
			copy(elem[size-lo:], tmp)
		} else {
			tmp := make([]E, size-lo)
			copy(tmp, elem[lo:size])
			copy(elem[size-lo:], elem[:hi])
			copy(elem, tmp)
		}

		q.head = size
		q.tail = 0
	}

	q.elem = elem
}

func (q *Queue[E]) index(i int) int {
	return i % len(q.elem)
}

func (q *Queue[E]) wraparound(i int) int {
	if i < 0 {
		i = -i
	}
	return i % (2 * len(q.elem))
}
