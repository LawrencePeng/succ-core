package util

import (
	"sort"
)

type QSufSort struct {
	I        []int32
	V        []int32
	Alphabet []int32
	R        int32
	H        int32
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

func (q *QSufSort) initAlphabet(set *HashSet) {
	q.Alphabet = make([]int32, set.Len())
	i := int32(0)
	for k := range set.M {
		q.Alphabet[i] = k
		i++
	}
	tr := make([]int, len(q.Alphabet))
	for i := 0; i < len(q.Alphabet); i++ {
		tr[i] = int(q.Alphabet[i])
	}
	sort.Ints(tr)
	for i := 0; i < len(q.Alphabet); i++ {
		q.Alphabet[i] = int32(tr[i])
	}
}

func (q *QSufSort) BuildSuffixArray(input Source) {
	max := int32(EOF)
	min := max

	q.I = make([]int32, input.Len()+2)
	q.V = make([]int32, input.Len()+2)

	alphabetSet := &HashSet{M: make(map[int32]bool)}
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
	if q.V[input.Len()] > int32(EOF) {
		max = q.V[input.Len()]
	}
	if q.V[input.Len()] < min {
		min = q.V[input.Len()]
	}
	alphabetSet.Add(int32(EOF))

	q.initAlphabet(alphabetSet)

	q.suffixSort(input.Len()+1, max+1, min)

}

func (q *QSufSort) SA() []int32 {
	return q.I
}

func (q *QSufSort) ISA() []int32 {
	return q.V
}

func (q *QSufSort) suffixSort(n int32, k int32, l int32) {
	var pi, pk, i, j, s, sl int32

	if n >= k-l {
		j = q.transform(n, k, l, n)
		q.bucketSort(n, j)
	} else {
		q.transform(n, k, l, 0x7fffffff)
		for i = 0; i < n; i++ {
			q.I[i] = i
		}
		q.H = 0
		q.sortSplit(0, n+1)
	}
	q.H = q.R

	for q.I[0] >= -n {
		pi = 0
		sl = 0

		s = q.I[pi]
		if s < 0 {
			pi -= s
			sl += s
		} else {
			if sl != 0 {
				q.I[pi+sl] = sl
				sl = 0
			}
			pk = q.V[s] + 1
			q.sortSplit(pi, pk-pi)
			pi = pk
		}
		for pi <= n {
			s = q.I[pi]
			if s < 0 {
				pi -= s
				sl += s
			} else {
				if sl != 0 {
					q.I[pi+sl] = sl
					sl = 0
				}
				pk = q.V[s] + 1
				q.sortSplit(pi, pk-pi)
				pi = pk
			}
		}

		if sl != 0 {
			q.I[pi+sl] = sl
		}
		q.H = 2 * q.H

	}

	for i = 0; i <= n; i++ {
		if q.V[i] > 0 {
			q.V[i]--
			q.I[q.V[i]] = i
		}
	}
}
func (q *QSufSort) sortSplit(p int32, n int32) {
	var pa, pb, pc, pd, pl, pm, pn int32
	var f, v, s, t int32

	if n < 7 {
		q.selectSortSplit(p, n)
		return
	}

	v = q.choosePivot(p, n)
	pb = p
	pa = pb

	pd = p + n - 1
	pc = pd

	for {
		if pb <= pc {
			f = q.Key(pb)
		}
		for pb <= pc && f <= v {
			if f == v {
				q.Swap(pa, pb)
				pa++
			}
			pb++
			if pb <= pc {
				f = q.Key(pb)
			}
		}
		if pc >= pb {
			f = q.Key(pc)
		}
		for pc >= pb && f >= v {
			if f == v {
				q.Swap(pc, pd)
				pd--
			}
			pc--
			if pb <= pc {
				f = q.Key(pc)
			}
		}
		if pb > pc {
			break
		}
		q.Swap(pb, pc)
		pb++
		pc--
	}
	pn = p + n

	s = pa - p
	t = pb - pa
	if s > t {
		s = t
	}

	pl = p
	pm = pb - s
	for s != 0 {
		q.Swap(pl, pm)
		s--
		pl++
		pm++
	}

	s = pd - pc
	t = pn - pd - 1
	if s > t {
		s = t
	}

	pl = pb
	pm = pn - s

	for s != 0 {
		q.Swap(pl, pm)
		s--
		pl++
		pm++
	}

	s = pb - pa
	t = pd - pc

	if s > 0 {
		q.sortSplit(p, s)
	}
	q.updateGroup(p+s, p+n-t-1)
	if t > 0 {
		q.sortSplit(p+n-t, t)
	}
}
func (q *QSufSort) choosePivot(p int32, n int32) int32 {
	var pl, pm, pn int32
	var s int32

	pm = p + int32(uint32(n)>>1)

	if n > 7 {
		pl = p
		pn = p + n - 1

		if n > 40 {
			s = int32(uint32(n) >> 3)
			pl = q.MED3(pl, pl+s, pl+s+s)
			pm = q.MED3(pm-s, pm, pm+s)
			pn = q.MED3(pn-s-s, pn-s, pn)
		}
		pm = q.MED3(pl, pm, pn)
	}
	return q.Key(pm)
}

func (q *QSufSort) MED3(a int32, b int32, c int32) int32 {
	if q.Key(a) < q.Key(b) {
		if q.Key(b) < q.Key(c) {
			return b
		} else {
			if q.Key(a) < q.Key(c) {
				return c
			} else {
				return a
			}
		}
	} else {
		if q.Key(b) > q.Key(c) {
			return b
		} else {
			if q.Key(a) > q.Key(c) {
				return c
			} else {
				return a
			}
		}
	}
}

func (q *QSufSort) selectSortSplit(p int32, n int32) {
	var pa, pb, pi, pn int32
	var f, v int32

	pa = p
	pn = p + n - 1

	for pa < pn {
		pb = pa + 1
		pi = pb
		for f = q.Key(pa); pi <= pn; pi++ {
			v = q.Key(pi)
			if v < f {
				f = v
				q.Swap(pi, pa)
				pb = pa + 1
			} else if v == f {
				q.Swap(pi, pb)
				pb++
			}
		}
		q.updateGroup(pa, pb-1)
		pa = pb
	}
	if pa == pn {
		q.V[q.I[pa]] = pa
		q.I[pa] = -1
	}
}
func (q *QSufSort) updateGroup(pl int32, pm int32) {
	var g int32

	g = pm
	q.V[q.I[pl]] = g
	if pl == pm {
		q.I[pl] = -1
	} else {
		pl++
		q.V[q.I[pl]] = g
		for pl < pm {
			pl++
			q.V[q.I[pl]] = g
		}
	}
}
func (q *QSufSort) Swap(a int32, b int32) {
	q.I[a], q.I[b] = q.I[b], q.I[a]
}
func (q *QSufSort) Key(p int32) int32 {
	a := q.I[p] + q.H
	b := q.V[a]
	return b
}

func (q *QSufSort) bucketSort(n int32, k int32) {
	var pi, i, c, d, g int32

	for pi = 0; pi < k; pi++ {
		q.I[pi] = -1
	}
	for i = 0; i <= n; i++ {
		c = q.V[i]
		q.V[i] = q.I[c]
		q.I[c] = i
	}
	pi = k - 1
	i = n
	for ; pi >= 0; pi-- {
		c = q.I[pi]
		d = q.V[c]
		g = i
		q.V[c] = g
		if d >= 0 {
			q.I[i] = c
			i--

			c = d
			d = q.V[c]
			q.V[c] = g
			q.I[i] = c
			i--
			for d >= 0 {
				c = d
				d = q.V[c]
				q.V[c] = g
				q.I[i] = c
				i--
			}
		} else {
			q.I[i] = -1
			i--
		}
	}
}

func (qS *QSufSort) transform(n int32, k int32, l int32, q int32) int32 {
	var b, c, d, e, i, j, m, s int32
	var pi, pj int32
	s = int32(0)
	for i = k - l; i != 0; i = i >> 1 {
		s++
	}
	e = int32(uint32(0x7fffffff) >> uint(s))
	qS.R = 0
	d = qS.R
	b = d
	c = d<<uint(s) | (k - l)

	for ; qS.R < n && d <= e && c <= q; qS.R++ {
		b = b<<uint(s) | (qS.V[qS.R] - l + 1)
		d = c

		if qS.R < n && d <= e {
			c = int32(uint32(d)<<uint(s)) | (k - l)
		}
	}

	m = (1 << uint((qS.R-1)*s)) - 1
	qS.V[n] = l - 1
	if d <= n {
		for pi = 0; pi <= d; pi++ {
			qS.I[pi] = 0
		}

		pi = qS.R
		c = b

		for ; pi <= n; pi++ {
			qS.I[c] = 1
			c = (c&m)<<uint(s) | (qS.V[pi] - l + 1)
		}

		for i = 1; i < qS.R; i++ {
			qS.I[c] = 1
			c = (c & m) << uint(s)
		}

		pi = 0
		j = 1
		for ; pi <= d; pi++ {
			if qS.I[pi] != 0 {
				qS.I[pi] = j
				j++
			}
		}

		pi = 0
		pj = qS.R
		c = b
		for pj <= n {
			qS.V[pi] = qS.I[c]
			c = (c&m)<<uint(s) | (qS.V[pj] - l + 1)
			pi++
			pj++
		}

		for pi < n {
			qS.V[pi] = qS.I[c]
			pi++
			c = (c & m) << uint(s)
		}

	} else {
		pi = 0
		pj = qS.R
		c = b

		for pj <= n {
			qS.V[pi] = c
			c = (c&m)<<uint(s) | (qS.V[pj] - l + 1)
			pi++
			pj++
		}

		for pi < n {
			qS.V[pi] = c
			pi++
			c = (c & m) << uint(s)
		}

		j = d + 1
	}
	qS.V[n] = 0
	return j
}
