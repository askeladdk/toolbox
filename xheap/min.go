package xheap

import "cmp"

// Min is a min-heap where elements are associated with priorities.
type Min[E any, P cmp.Ordered] []Pair[E, P]

// Empty reports whether h is empty.
func (h Min[E, P]) Empty() bool {
	return len(h) == 0
}

// Fix re-establishes the heap ordering after the priority of the i-th element has changed.
func (h Min[E, P]) Fix(i int) {
	Fix(h.iface(), i)
}

// Init establishes the min-heap invariants.
func (h Min[E, P]) Init() {
	Init(h.iface())
}

// Len reports the number of elements in h.
func (h Min[E, P]) Len() int {
	return len(h)
}

// Peek returns the minimum element without removing it.
// Panics if the heap is empty.
func (h Min[E, P]) Peek() (value E, priority P) {
	return h[0].Value, h[0].Priority
}

// Pop removes and returns the minimum element.
// Panics if the heap is empty.
func (h *Min[E, P]) Pop() (value E, priority P) {
	x := Pop(h.iface())
	return x.Value, x.Priority
}

// Push pushes a new element on the heap.
func (h *Min[E, P]) Push(value E, priority P) {
	Push(h.iface(), Pair[E, P]{value, priority})
}

// Remove removes and returns the i-th element.
// Panics if the heap is empty.
func (h *Min[E, P]) Remove(i int) (value E, priority P) {
	p := Remove(h.iface(), i)
	return p.Value, p.Priority
}

// Reset empties the heap.
func (h *Min[E, P]) Reset() {
	*h = (*h)[:0]
}

func (h *Min[E, P]) iface() Interface[Pair[E, P]] {
	return (*minHeapImpl[E, P])(h)
}

type minHeapImpl[E any, P cmp.Ordered] []Pair[E, P]

func (h minHeapImpl[E, P]) Len() int {
	return len(h)
}

func (h minHeapImpl[E, P]) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h *minHeapImpl[E, P]) Pop() (x Pair[E, P]) {
	n := len(*h) - 1
	*h, x = (*h)[:n], (*h)[n]
	return x
}

func (h *minHeapImpl[E, P]) Push(x Pair[E, P]) {
	*h = append(*h, x)
}

func (h minHeapImpl[E, P]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
