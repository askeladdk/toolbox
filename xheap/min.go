package xheap

import "cmp"

// Min is a min-heap where an element is associated with a priority.
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
func (h Min[E, P]) Peek() (E, P) {
	return h[0].Elem, h[0].Prio
}

// Pop removes and returns the minimum element.
func (h *Min[E, P]) Pop() (E, P) {
	x := Pop(h.iface())
	return x.Elem, x.Prio
}

// Push pushes a new element on the heap.
func (h *Min[E, P]) Push(elem E, prio P) {
	Push(h.iface(), Pair[E, P]{elem, prio})
}

// Remove removes and returns the i-th element.
func (h *Min[E, P]) Remove(i int) (E, P) {
	p := Remove(h.iface(), i)
	return p.Elem, p.Prio
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
	return h[i].Prio < h[j].Prio
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
