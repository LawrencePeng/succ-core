package util_test
import (
	"testing"
	"bytes"
	"github.com/lexandro/go-assert"
	"."
)

var testSizeInBits int64 = 1024 * 1024

func TestBitVecGetBit(t *testing.T) {
	bv := util.NewBitVectorWithSize(int64(testSizeInBits))
	for i := int64(0); i < testSizeInBits; i++ {
		if i % 2 == 0 {
			bv.SetBit(i)
		}
	}

	buf := new(bytes.Buffer)
	bv.WriteToBuf(buf)

	bufBlock := util.ReadInt(buf)

	data := util.ToLongSlice(buf.Next(int(bufBlock) * util.LONG_SIZE))

	for i := int64(0); i < testSizeInBits; i++ {
		val := util.BitVecGetBit(data, i)
		var expectedVal int64

		if i % 2 == 0 {
			expectedVal = 1
		} else {
			expectedVal = 0
		}

		assert.Equals(t, expectedVal, val)
	}
}

func TestNumBlocks2(t *testing.T) {
	bv := util.NewBitVectorWithSize(42)
	assert.Equals(t, len(bv.Data), 1)
}


func TestBitVecGetValue(t *testing.T) {
	pos := int64(0)
	bv := util.NewBitVectorWithSize(testSizeInBits)
	for i := 0; i < 10000; i++ {
		bv.SetValue(pos, int64(i), util.BitWidth(int64(i)))
		pos += int64(util.BitWidth(int64(i)))
	}

	pos = 0

	for i := int32(0); i < int32(10000); i++ {
		val := bv.GetValue(pos, util.BitWidth(int64(i)))
		assert.Equals(t, int64(i), val)
		pos += int64(util.BitWidth(int64(i)))
	}
}

