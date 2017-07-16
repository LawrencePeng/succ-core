package util_test

import "testing"
import (
	"github.com/lexandro/go-assert"
)
import "."

func TestGetAndSetBit(t *testing.T) {
	n := int64(0)

	for i := 0; i < 64; i++ {
		assert.Equals(t, 0, int(util.GetBit(n , int32(i))))
	}

	for i := 0; i < 64; i++ {
		n = util.SetBit(n, int32(i))
		for j := 0; j < 64; j++ {
			if i == j {
				assert.Equals(t, 1, int(util.GetBit(n, int32(j))))
			} else {
				assert.Equals(t, 0, int(util.GetBit(n, int32(j))))
			}
		}
		n = 0
	}

	for i := 0; i < 64; i++ {
		n = util.SetBit(n, int32(i))
		for j := 0; j < 64; j++ {
			if j <= i {
				assert.Equals(t, 1, int(util.GetBit(n, int32(j))))
			} else {
				assert.Equals(t, 0, int(util.GetBit(n, int32(j))))
			}
		}
	}
}

func TestBitsToBlock64(t *testing.T) {
	assert.Equals(t,0, int(util.BitsToBlock64(0)))
	for i := 1; i <= 64; i++ {
		assert.Equals(t,1, int(util.BitsToBlock64(int64(i))))
	}
	assert.Equals(t,2, int(util.BitsToBlock64(int64(65))))
	assert.Equals(t, 4, int(util.BitsToBlock64(int64(256))))
	assert.Equals(t, 5, int(util.BitsToBlock64(int64(257))))
}

func TestBitWidth(t *testing.T) {
	assert.Equals(t,int32(1), util.BitWidth(0))
	assert.Equals(t,int32(1), util.BitWidth(1))
	assert.Equals(t,int32(2), util.BitWidth(2))
	assert.Equals(t,int32(2), util.BitWidth(3))
	assert.Equals(t,int32(3), util.BitWidth(4))
	assert.Equals(t,int32(8), util.BitWidth(255))
	assert.Equals(t, int32(9), util.BitWidth(256))
}