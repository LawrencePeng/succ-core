package util

type BMArray struct {
	BM   *BitMap
	bits int
	n    int
}

func NewBMArray(n, bits int) *BMArray {
	return &BMArray{
		BM:   NewBitMap(int64(n * bits)),
		bits: bits,
		n:    n,
	}
}

func NewBMAArrayNative(input []int, offset, length int) *BMArray {
	ba := &BMArray{
		BM:   NewBitMap(int64(length * IntLog2(int64(length)))),
		n:    length,
		bits: IntLog2(int64(length + 1)),
	}
	for i := offset; i < offset+length; i++ {
		ba.SetVal(i-offset, int64(input[i]))
	}
	return ba
}
func (array *BMArray) SetVal(i int, val int64) {
	s := i * array.bits
	e := s + array.bits - 1

	if s/64 == e/64 {
		array.BM.Data[s/64] |= val << uint64(63-e%64)
	} else {
		array.BM.Data[s/64] |= int64(uint64(val) >> uint64(63-e%64))
		array.BM.Data[e/64] |= val << uint64(63-e%64)
	}
}

func (array *BMArray) GetVal(i int) int64 {
	var val int64
	s := i * array.bits
	e := s + array.bits - 1

	if s/64 == e/64 {
		val = array.BM.Data[s/64] << uint64(s%64)
		val = int64(uint64(val) >> uint64(63-e%64+s%64))
	} else {
		val1 := array.BM.Data[s/64] << uint64(s%64)
		val2 := int64(uint64(array.BM.Data[e/64]) >> uint64(63-e%64))
		val1 = int64(uint64(val1 >> uint64(s%64-(e%64+1))))
		val = val1 | val2
	}

	return val
}
