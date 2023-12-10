package xheap

import "cmp"

// Pair associated a value with a priority used to order the heap.
type Pair[E any, P cmp.Ordered] struct {
	Elem E
	Prio P
}

// Max is a max-heap where an element is associated with a priority.
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
func (h Max[E, P]) Peek() (E, P) {
	return h[0].Elem, h[0].Prio
}

// Pop removes and returns the maximum element.
func (h *Max[E, P]) Pop() (E, P) {
	x := Pop(h.iface())
	return x.Elem, x.Prio
}

// Push pushes a new element on the heap.
func (h *Max[E, P]) Push(elem E, prio P) {
	Push(h.iface(), Pair[E, P]{elem, prio})
}

// Remove removes and returns the i-th element.
func (h *Max[E, P]) Remove(i int) (E, P) {
	p := Remove(h.iface(), i)
	return p.Elem, p.Prio
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
	return h[j].Prio < h[i].Prio
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
