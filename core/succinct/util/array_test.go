package util

import (
	"testing"
	"github.com/lexandro/go-assert"
)

func GetRank132(arrayBuf []int32, startPos int32, size int32, i int32) int32 {
	sp := int32(0)
	ep := size - 1

	var m int32

	for ; sp <= ep; {
		m = (sp + ep) / 2
		if arrayBuf[startPos + m] == i {
			return m + 1
		} else if i < arrayBuf[startPos + m] {
			ep = m - 1
		} else {
			sp = m + 1
		}
	}

	return ep + 1
}

func TestGetRank132(t *testing.T) {
	data := []int32{2,3,5,7,11,13,17,19,23,29}
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(0)), int32(0))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(2)), int32(1))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(3)), int32(2))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(4)), int32(2))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(6)), int32(3))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(22)), int32(8))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(29)), int32(10))
	assert.Equals(t, GetRank132(data, 0, int32(len(data)), int32(33)), int32(10))

}