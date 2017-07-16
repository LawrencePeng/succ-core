package util

import "bytes"

type IntVector struct {
	Bv			*BitVector
	BitWidth	int32
}

func NewIntVector(numElements, bitWidth int32) *IntVector {
	bv := NewBitVectorWithSize(int64(numElements * bitWidth))
	return &IntVector{
		Bv: bv,
		BitWidth: bitWidth,
	}
}

func CopyFromBitVector(bitWidth int32, bitVector *BitVector) *IntVector {
	bv := CopyBitVector(bitVector)
	return &IntVector{
		Bv:bv,
		BitWidth:bitWidth,
	}
}

func (iv *IntVector) Add(idx, ele int32) {
	iv.Bv.SetValue(int64(idx * iv.BitWidth), int64(ele), int32(iv.BitWidth))
}

func CopyFromIntArrayList(data *IntArrayList, bitWidth int32) *IntVector {
	bv := NewBitVectorWithSize(int64(data.Size() * bitWidth))
	iv := &IntVector{
		Bv:       bv,
		BitWidth: bitWidth,
	}

	for i := int32(0); i < data.Size(); i++ {
		if BitWidth(int64(data.Get(i))) > bitWidth {
			panic("Ill Format Precision")
		}
		iv.Add(i, data.Get(i))
	}
	return iv
}

func ReadIntVectorFromBuf(buf *bytes.Buffer) *IntVector {
	bitWidth := ReadInt(buf)
	if bitWidth == 0 {
		return nil
	}
	return CopyFromBitVector(bitWidth, ReadBitVectorFromBuf(buf))
}

func (iv *IntVector) Get(idx int32) int32 {
	return int32(iv.Bv.GetValue(int64(idx * iv.BitWidth), iv.BitWidth))
}

func (iv *IntVector) SerializedSize() int32 {
	return int32(INT_SIZE) + iv.Bv.SerializedSize()
}

func (iv *IntVector) WriteToBuf(buf *bytes.Buffer) {
	WriteInt(buf, iv.BitWidth)
	iv.Bv.WriteToBuf(buf)
}

func IntVecGet(data []int64, index int32, bitWidth int32) int32 {
	return int32(BitVecGetValue(data, int64(index * bitWidth), bitWidth))
}