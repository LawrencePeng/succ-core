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

func (bv *BitVector) SetBit(i int64) {
	bv.Data[i >> 6] = SetBit(bv.Data[int64(uint64(i) >> 6)], int32(i % 64))
}

func (bv *BitVector) ClearBit(i int64) {
	bv.Data[i >> 6] = ClearBit(bv.Data[int64(uint64(i) >> 6)], int(i % 64))
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
		bv.Data[sIdx] = (bv.Data[sIdx] & (LOW_BITS_SET[sOff] | LOW_BITS_UNSET[sOff + int64(bits)])) |
			value << uint64(sOff)
	} else {
		bv.Data[sIdx] = (bv.Data[sIdx] & LOW_BITS_SET[sOff]) | value << uint64(sOff)
		bv.Data[sIdx + 1] =
			(bv.Data[sIdx + 1] & LOW_BITS_UNSET[(sOff + int64(bits)) % 64]) | int64(uint64(value) >> (64 - uint(sOff)))
	}
}

func (bv *BitVector) GetValue(pos int64, bits int32) int64 {
	if bits == 0 {
		return 0
	}

	sOff := int(pos % 64)
	sIdx := int(pos / 64)

	if sOff + sIdx <= 64 {
		return int64(uint64(bv.Data[sIdx]) >> uint64(sOff)) & LOW_BITS_SET[bits]
	}
	return int64(uint64(bv.Data[sIdx] >> uint(sOff))) |
		bv.Data[sIdx + 1] << uint(64 - sOff) &
		LOW_BITS_UNSET[bits]
}

func (bv *BitVector) SerializedSize() int32 {
	return int32(8 * len(bv.Data) +  int(LONG_SIZE))
}

func BitVecGetBit(data []int64, i int64) int64 {
	return GetBit(data[int32(uint64(i) >> 6)], int32(i % 64))
}

func BitVecGetValue(data []int64, pos int64, bits int32) int64 {
	if bits == 0 {
		return 0
	}

	sOff := int32(pos % 64)
	sIdx := int32(pos / 64)

	if sOff + sIdx <= 64 {
		return data[int32(uint32(sIdx) >> uint(sOff))] & LOW_BITS_SET[bits]
	}
	return data[int32(uint32(sIdx) >> uint(sOff))] | data[sIdx + 1] << uint(64 - sOff)
}