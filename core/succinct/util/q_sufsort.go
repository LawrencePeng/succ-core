package util

type QSufSort struct {
	I 			[]int32
	V 			[]int32
	alphabet	[]int32
	r			int32
	h 			int32
}

type HashSet struct {
	m map[int32]bool
}

func (h *HashSet) Len() int {
	return len(h.m)
}

func (q *QSufSort) initAlphabet(set HashSet) {
	q.alphabet = make([]int32, set.Len())
	i := int32(0)
	for
}

