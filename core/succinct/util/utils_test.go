package util

import "testing"
import (
	"github.com/lexandro/go-assert"
)

func TestGetAndSetBit(t *testing.T) {
	n := int64(0)

	for i := 0; i < 64; i++ {
		assert.Equals(t, 0, int(GetBit(n , i)))
	}

	for i := 0; i < 64; i++ {
		n = SetBit(n, i)
		for j := 0; j < 64; j++ {
			if i == j {
				assert.Equals(t, 1, int(GetBit(n, j)))
			} else {
				assert.Equals(t, 0, int(GetBit(n, j)))
			}
		}
		n = 0
	}

	for i := 0; i < 64; i++ {
		n = SetBit(n, i)
		for j := 0; j < 64; j++ {
			if j <= i {
				assert.Equals(t, 1, int(GetBit(n, j)))
			} else {
				assert.Equals(t, 0, int(GetBit(n, j)))
			}
		}
	}
}

func TestBitsToBlock64(t *testing.T) {
	assert.Equals(t,0, int(BitsToBlock64(0)))
	for i := 1; i <= 64; i++ {
		assert.Equals(t,1, int(BitsToBlock64(int64(i))))
	}
	assert.Equals(t,2, int(BitsToBlock64(int64(65))))
	assert.Equals(t, 4, int(BitsToBlock64(int64(256))))
	assert.Equals(t, 5, int(BitsToBlock64(int64(257))))
}

func TestBitWidth(t *testing.T) {
	assert.Equals(t,1, BitWidth(0))
	assert.Equals(t,1, BitWidth(1))
	assert.Equals(t,2, BitWidth(2))
	assert.Equals(t,2, BitWidth(3))
	assert.Equals(t,3, BitWidth(4))
	assert.Equals(t,8, BitWidth(255))
	assert.Equals(t, 9, BitWidth(256))
}