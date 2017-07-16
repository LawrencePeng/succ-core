package util_test

import "testing"
import (
	. "."
	"github.com/lexandro/go-assert"
)

func TestList(t *testing.T) {
	list := NewIntArrayList()
	for i := 0; i < 1000000; i++ {
		list.Add(int32(i))
	}
	assert.Equals(t, int32(1000000), int32(list.Size()))
	for i := 0; i < 1000000; i++ {
		assert.Equals(t, int32(i), list.Get(int32(i)))
	}
}
