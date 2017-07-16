package util

import (
	"bytes"
)

type BitVector struct {
	Data []int64
}

func NewBitVectorWithSize(sizeInBits int64) *BitVector {
	bv := &BitVector{}
	bv.Data = make([]int64, BitsToBlock64(sizeInBits))
	return bv
}

func CopyBitVector(ano *BitVector) *BitVector {
	dst_data := make([]int64, len(ano.Data))
	bv := &BitVector{dst_data}
	return bv
}

func ReadBitVectorFromBuf(buf *bytes.Buffer) *BitVector {
	numOfBlock := ReadInt(buf)

	if numOfBlock == 0 {
		return nil
	}

	data := make([]int64, numOfBlock)
	for i := 0; i < len(data); i++ {
		data[i] = ReadLong(buf)
	}

	return &BitVector {
		Data: data,
	}
}

func (bv *BitVector) WriteToBuf(buf *bytes.Buffer) {
	WriteInt(buf, int32(len(bv.Data)))
	for _, n := range bv.Data {
		WriteLong(buf, n)
	}
}

func (bv *BitVector) WriteToDataBuf(buf *bytes.Buffer) {
	for _, n := range bv.Data {
		WriteLong(buf, n)
	}
}

func (bv *BitVector) SetBit(i int64) {
	bv.Data[i >> 6] = SetBit(bv.Data[int64(uint64(i) >> 6)], int32(i % 64))
}

func (bv *BitVector) ClearBit(i int64) {
	bv.Data[i >> 6] = ClearBit(bv.Data[int64(uint64(i) >> 6)], int32(i % 64))
}

func (bv *BitVector) GetBit(i int64) int64 {
	return GetBit(bv.Data[int64(uint64(i) >> 6)], int32(i % 64))
}

func (bv *BitVector) SetValue(pos, value int64, bits int32) {
	if bits == 0 {
		return
	}

	sOff := int64(pos % 64)
	sIdx := int64(pos / 64)

	if sOff + int64(bits) <= 64 {
		a := bv.Data[sIdx]
		b := int64(LOW_BITS_SET[sOff])
		c := int64(LOW_BITS_UNSET[sOff + int64(bits)])
		d := int64(value) << uint64(sOff)
		e := int64((a & (b | c)) | d)
		bv.Data[sIdx] = e

	} else {
		a := int64(bv.Data[sIdx])
		b := int64(LOW_BITS_SET[sOff])
		c := int64(value) << uint(sOff)
		d := (a & b) | c
		bv.Data[sIdx] = int64(d)

		e := int64(bv.Data[sIdx + 1])
		f := int64(LOW_BITS_UNSET[(sOff + int64(bits)) % 64])
		g := int64(value) >> uint(64 - sOff)
		bv.Data[sIdx + 1] = int64((e & f) | g)
	}
}

func (bv *BitVector) GetValue(pos int64, bits int32) int64 {
	if bits == 0 {
		return 0
	}

	sOff := int(pos % 64)
	sIdx := int(pos / 64)

	if sOff + int(bits) <= 64 {
		a := int64(uint64(bv.Data[sIdx]) >> uint(sOff))
		b := int64(LOW_BITS_SET[bits])
		c := int64(a & b)
		return c
	}

	a := int64(uint64(bv.Data[sIdx]) >> uint(sOff))
	b := int64(bv.Data[sIdx + 1] << uint(64 - sOff))
	c := int64(LOW_BITS_SET[bits])
	d := int64((a | b) & c)

	return d

}

func (bv *BitVector) SerializedSize() int32 {
	return int32(8 * len(bv.Data) + int(INT_SIZE))
}

func BitVecGetBit(data []int64, i int64) int64 {
	return GetBit(data[int64(uint64(i) >> 6)], int32(i % 64))
}

func BitVecGetValue(data []int64, pos int64, bits int32) int64 {
	if bits == 0 {
		return 0
	}

	sOff := int(pos % 64)
	sIdx := int(pos / 64)

	if sOff + int(bits) <= 64 {
		a := int64(uint64(data[sIdx]) >> uint(sOff))
		b := int64(LOW_BITS_SET[bits])
		c := int64(a & b)
		return c
	}

	a := int64(uint64(data[sIdx]) >> uint(sOff))
	b := int64(data[sIdx + 1] << uint(64 - sOff))
	c := int64(LOW_BITS_SET[bits])
	d := int64((a | b) & c)

	return d
}