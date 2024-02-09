package xheap

import "cmp"

// Pair associated a value with a priority used to order the heap.
type Pair[E any, P cmp.Ordered] struct {
	Value    E
	Priority P
}

// Max is a max-heap where values are associated with a priorities.
type Max[E any, P cmp.Ordered] []Pair[E, P]

// Empty reports whether h is empty.
func (h Max[E, P]) Empty() bool {
	return len(h) == 0
}

// Fix re-establishes the heap ordering after the priority of the i-th element has changed.
func (h Max[E, P]) Fix(i int) {
	Fix(h.iface(), i)
}

// Init establishes the max-heap invariants.
func (h Max[E, P]) Init() {
	Init(h.iface())
}

// Len reports the number of elements in h.
func (h Max[E, P]) Len() int {
	return len(h)
}

// Peek returns the maximum element without removing it.
// Panics if the heap is empty.
func (h Max[E, P]) Peek() (value E, priority P) {
	return h[0].Value, h[0].Priority
}

// Pop removes and returns the maximum element.
// Panics if the heap is empty.
func (h *Max[E, P]) Pop() (value E, priority P) {
	x := Pop(h.iface())
	return x.Value, x.Priority
}

// Push pushes a new element on the heap.
func (h *Max[E, P]) Push(value E, priority P) {
	Push(h.iface(), Pair[E, P]{value, priority})
}

// Remove removes and returns the i-th element.
// Panics if the heap is empty.
func (h *Max[E, P]) Remove(i int) (value E, priority P) {
	p := Remove(h.iface(), i)
	return p.Value, p.Priority
}

// Reset empties the heap.
func (h *Max[E, P]) Reset() {
	*h = (*h)[:0]
}

func (h *Max[E, P]) iface() Interface[Pair[E, P]] {
	return (*maxHeapImpl[E, P])(h)
}

type maxHeapImpl[E any, P cmp.Ordered] []Pair[E, P]

func (h maxHeapImpl[E, P]) Len() int {
	return len(h)
}

func (h maxHeapImpl[E, P]) Less(i, j int) bool {
	return h[j].Priority < h[i].Priority
}

func (h *maxHeapImpl[E, P]) Pop() (x Pair[E, P]) {
	n := len(*h) - 1
	*h, x = (*h)[:n], (*h)[n]
	return x
}

func (h *maxHeapImpl[E, P]) Push(x Pair[E, P]) {
	*h = append(*h, x)
}

func (h maxHeapImpl[E, P]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
