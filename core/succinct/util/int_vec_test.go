package util_test

import (
	. "."
	"testing"
	"github.com/lexandro/go-assert"
	"bytes"
)

const (
	testSize = int32(1024 * 1024)
	testBits = int32(20)
)



func TestIntVecGet(t *testing.T) {
	iv := NewIntVector(testSize, testBits)

	for i := int32(0); i < testSize ; i++ {
		iv.Add(i, i)
	}

	buf := new(bytes.Buffer)
	iv.WriteToBuf(buf)

	bitWidth := ReadInt(buf)
	bufBlocks := ReadInt(buf)

	data := ToLongSlice(buf.Next(int(bufBlocks) * LONG_SIZE))
	for i := int32(0); i < int32(100); i++ {
		val := IntVecGet(data, i, bitWidth)
		assert.Equals(t, i, val)
	}
}

func TestIntVectorAddAndGet(t *testing.T) {
	iv := NewIntVector(testSize, testBits)

	for i := int32(0); i < testSize; i++ {
		iv.Add(i, i)
	}

	for i := int32(0); i < 10000; i++ {
		if iv.Get(i) != i {
			assert.Equals(t, i, iv.Get(i))
		}
	}
}