package util

import "bytes"

type BitMap struct {
	Data []int64
	Size int64
}

func NewBitMap(size int64) *BitMap {
	bmSize := size/64 + 1
	return &BitMap{
		Data: make([]int64, bmSize),
		Size: size,
	}
}

func (bm *BitMap) BitMapSize() int {
	return len(bm.Data)
}

func (bm *BitMap) SetBit(i int) {
	bm.Data[i/64] |= 1 << uint(63-i)
}

func (bm *BitMap) GetBit(i int) int64 {
	return int64(uint64(bm.Data[i/64])>>uint(63-i)) & 1
}

func (bm *BitMap) SetValPos(pos int, val int64, bits int) {
	e := int64(pos) + int64(bits) - 1
	if int64(pos)/64 == e/64 {
		bm.Data[pos/64] |= val << uint64(63-e&64)
	} else {
		bm.Data[pos/64] |= int64(uint64(val) >> uint64(e%64+1))
		bm.Data[e/64] |= val << uint64(63-e%64)
	}
}

func (bm *BitMap) GetValPos(pos, bits int) int64 {
	var val int64
	s := int64(pos)
	e := s + int64(bits-1)

	val = bm.Data[s/64] << uint64(s%64)
	if s/64 == e/64 {
		val = int64(uint64(val) >> uint64(63-e%64+s%64))
	} else {
		val = int64(uint64(val) >> uint64(s%64-(e%64+1)))
		val |= int64(uint64(bm.Data[e/64]) >> uint64(63-e%64))
	}
	return val
}

func (bm *BitMap) GetSelect1(i int) int64 {
	sel, count := int64(-1), int64(0)
	for j := 0; j < int(bm.Size); j++ {
		if bm.GetBit(j) == 1 {
			count++
		}
		if int(count) == i+1 {
			sel = int64(j)
			break
		}
	}
	return sel
}

func (bm *BitMap) GetSelect0(i int) int64 {
	sel, count := int64(-1), int64(0)
	for j := 0; j < int(bm.Size); j++ {
		if bm.GetBit(j) == 0 {
			count++
		}
		if int(count) == i+1 {
			sel = int64(j)
			break
		}
	}
	return sel
}

func (bm *BitMap) GetRank1(i int) int64 {
	count := int64(0)

	for j := 0; j <= i; j++ {
		if bm.GetBit(j) == 1 {
			count++
		}
	}

	return count
}

func (bm *BitMap) GetRank0(i int) int64 {
	count := int64(0)

	for j := 0; j <= i; j++ {
		if bm.GetBit(j) == 0 {
			count++
		}
	}

	return count
}

func (bm *BitMap) Clear() {
	for i := 0; i < len(bm.Data); i++ {
		bm.Data[i] = 0
	}
}

type LongBuffer struct {
	bytes.Buffer
	len int
}

func (lb *LongBuffer) Len() int {
	return lb.len
}
