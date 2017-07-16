package util_test

import "testing"

import (
	. "."
	"github.com/lexandro/go-assert"
	"bytes"
	"fmt"
	"time"
)

const deltestSize = int32(1024 * 1024)

func TestDeltaIntVector_Get(t *testing.T) {
	data := make([]int32, deltestSize)
	for i := int32(0); i < deltestSize; i++ {
		data[i] = i
	}

	vec := NewDeltaIntVector(data, 128)

	for i := int32(0); i < deltestSize; i++ {
		assert.Equals(t, i, vec.Get(i))
	}

	vec2 := NewDeltaIntVectorFull(data, 256 * 1024, 512 * 1024, 128)

	for i := int32(256 * 1024); i < 768*1024; i++ {
		assert.Equals(t, i, vec2.Get(i - 256 * 1024))
	}
}

func TestDeltaGet(t *testing.T) {
	data := make([]int32, deltestSize)
	for i := int32(0); i < deltestSize; i++ {
		data[i] = i
	}

	vec := NewDeltaIntVector(data, 128)

	buf := new(bytes.Buffer)

	vec.WriteToBuf(buf)
	bts := buf.Bytes()
	tic := time.Now()
	for i := int32(0); i < 1024; i++ {
		assert.Equals(t, i, DIVGet(&bts, i))
	}
	toc := time.Now()
	fmt.Print(toc.Sub(tic))

}

func TestBinarySearch(t *testing.T) {
	data := make([]int32, deltestSize)
	for i := int32(0); i < deltestSize; i++ {
		data[i] = i
	}

	vec := NewDeltaIntVector(data, 128)

	buf := new(bytes.Buffer)
	vec.WriteToBuf(buf)
	ptr := buf.Bytes()
	for i := int32(0); i < deltestSize; i ++ {
		val := BinarySearch(&ptr, i, 0, deltestSize - 1, true)
		assert.Equals(t, i, val)
	}

	for i := int32(0); i < deltestSize; i++ {
		data[i] = i * 2
	}

	vec = NewDeltaIntVector(data, 128)

	buf = new(bytes.Buffer)
	vec.WriteToBuf(buf)
	ptr = buf.Bytes()

	for i := int32(0); i < deltestSize; i++ {
		lo := BinarySearch(&ptr, 2 * i + 1, 0, deltestSize - 1, true)
		hi := BinarySearch(&ptr, 2 * i + 1, 0, deltestSize - 1, false)
		assert.Equals(t, i, lo)
		assert.Equals(t, i + 1, hi)
	}
}