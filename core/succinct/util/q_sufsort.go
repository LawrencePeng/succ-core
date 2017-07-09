package util

import "sort"

type QSufSort struct {
	I 			[]int32
	V 			[]int32
	Alphabet	[]int32
	R			int32
	H 			int32
}

type HashSet struct {
	M map[int32]bool
}

func (h *HashSet) Len() int {
	return len(h.M)
}

func (h *HashSet) Add(i int32) {
	h.M[i] = true
}

func (q *QSufSort) initAlphabet(set HashSet) {
	q.Alphabet = make([]int32, set.Len())
	i := int32(0)
	for k := range set.M {
		q.Alphabet[i] = k
		i++
	}
	sort.Sort(q.Alphabet)
}

func (q *QSufSort) BuildSuffixArray(input Source) {
	max := int32(EOF)
	min := max

	q.I = make([]int32, input.Len(), + 2)
	q.V = make([]int32, input.Len() + 2)

	alphabetSet := &HashSet{}
	for i := int32(0); i < input.Len(); i++ {
		q.V[i] = input.Get(i)
		if q.V[i] > max {
			max = q.V[i]
		}
		if q.V[i] < min {
			min = q.V[i]
		}
		alphabetSet.Add(input.Get(i))
	}
	q.V[input.Len()] = int32(EOF)


}

